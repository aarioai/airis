package debug

import (
	"bufio"
	"fmt"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/ios"
	"github.com/aarioai/airis/pkg/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	defaultLogRetention           = -time.Hour * 24 * 90
	defaultLogFileDateFormat      = "2006-01-02"
	defaultLogSuffix              = ".log"
	defaultLogBufferFlushInterval = time.Second * 3
)

type LogWriter struct {
	retention     time.Duration
	perm          os.FileMode
	dir           string
	symlink       string
	file          *os.File
	logName       string
	mu            sync.RWMutex
	cleanupTicker *time.Ticker
	done          chan struct{}
	bufSize       int
	bufWriter     *bufio.Writer
	lastFlush     time.Time
	flushInterval time.Duration

	dateFormat string
	suffix     string
}

var (
	onceLogWriter sync.Once
	logWriter     *LogWriter
)

func Console(args ...any) {
	if len(args) == 1 {
		switch v := args[0].(type) {
		case error:
			Console(v.Error())
			return
		case string:
			if v == "" {
				return
			}
		case nil:
			return
		}
	}

	msg := afmt.SprintfArgs(args)
	if msg == "" {
		return
	}

	msg = strings.TrimSuffix(msg, "\n")

	// 方便运行程序时直接显示
	fmt.Println(msg)
	// 记录进日志，方便通过消息队列通知
	log.Println(msg)
}

func RedirectLog(dir string, perm os.FileMode, bufSize int, symlink ...string) error {
	var err error
	if dir, err = ios.MkdirAll(dir, perm); err != nil {
		return err
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".log") {
			return nil
		}
		if info.Size() == 0 {
			os.Remove(path)
		}
		return nil
	})

	panicLog := path.Join(dir, fmt.Sprintf("panic-%s.log", time.Now().Format("2006-01-02")))
	var panicFile *os.File
	if panicFile, err = os.OpenFile(panicLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, perm); err != nil {
		return fmt.Errorf("failed to open crash file: %w", err)
	}
	defer panicFile.Close()

	if err = redirectStderr(panicFile); err != nil {
		return fmt.Errorf("failed to redirect stderr: %w", err)
	}
	var lw *LogWriter
	if lw, err = NewLogWriter(dir, perm, bufSize, symlink...); err != nil {
		return fmt.Errorf("failed to create log writer: %w", err)
	}
	log.SetOutput(lw)
	log.SetFlags(log.Ltime)
	return nil
}

// NewLogWriter 单例模式创建新的日志写入器
// @warn docker 内使用，需要注意单独添加软链接映射，如 -v /var/log/symlink.log:/var/log/symlink.log
func NewLogWriter(dir string, perm os.FileMode, bufSize int, symlinks ...string) (*LogWriter, error) {
	var err error
	defer func() {
		if err = generateRmlogScript(dir); err != nil {
			Console("generate rmlog script failed: %v", err)
		}
	}()

	if dir, err = ios.MkdirAll(dir, perm); err != nil {
		return nil, err
	}
	// 这里会同时判断 symlinks[0] 是否为空字符串，兼容性更强
	symlink := afmt.First(symlinks)
	if symlink == "" {
		symlink = path.Join(dir, "app.log")
	} else {
		symlink, _, err = ios.PrepareFile(symlink, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	onceLogWriter.Do(func() {
		logWriter = &LogWriter{
			retention:     defaultLogRetention,
			perm:          perm,
			dir:           dir,
			symlink:       symlink,
			bufSize:       bufSize, // 这个可以决定输出方式，因此尽量实例化之初指定，不要后来重新指定
			done:          make(chan struct{}),
			flushInterval: defaultLogBufferFlushInterval,

			dateFormat: defaultLogFileDateFormat,
			suffix:     defaultLogSuffix,
		}
		logWriter.initialize()
		logWriter.startCleanupRoutine()
	})

	if needsUpdate := logWriter.updateConfig(dir, symlink, bufSize); needsUpdate {
		logWriter.initialize()
	}

	return logWriter, nil
}
func (lw *LogWriter) WithRetentionDays(retentionDays time.Duration) *LogWriter {
	lw.retention = retentionDays
	return lw
}
func (lw *LogWriter) WithFlushInterval(flushInterval time.Duration) *LogWriter {
	lw.flushInterval = flushInterval
	return lw
}
func (lw *LogWriter) WithDateFormat(dateFormat string) *LogWriter {
	lw.dateFormat = dateFormat
	return lw
}
func (lw *LogWriter) WithSuffix(suffix string) *LogWriter {
	lw.suffix = suffix
	return lw
}
func (lw *LogWriter) initialize() {
	if err := lw.openFile(); err != nil {
		Console("initialize log file: %v", err)
	}
}

func (lw *LogWriter) updateConfig(dir, symlink string, bufSize int) bool {
	if lw.dir == dir && lw.symlink == symlink && lw.bufSize == bufSize {
		return false
	}

	Console("update log writer config: dir=%s, symlink=%s, bufSize=%d", dir, symlink, bufSize)
	lw.dir = dir
	lw.symlink = symlink
	lw.bufSize = bufSize
	return true
}

func (lw *LogWriter) currentPath() string {
	return path.Join(lw.dir, time.Now().Format(lw.dateFormat)+lw.suffix)
}
func (lw *LogWriter) write(p []byte) (n int, err error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	newLogFile := lw.currentPath()
	if newLogFile != lw.logName || lw.file == nil {
		if err = lw.openFile(); err != nil {
			Console("log writer failed to openFile log: %v", err)
		}
	}
	return lw.file.Write(p)
}
func (lw *LogWriter) writeBuffer(p []byte) (n int, err error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.currentPath() != lw.logName || lw.file == nil || lw.bufWriter == nil {
		if err := lw.openFile(); err != nil {
			return 0, fmt.Errorf("log writer failed to openFile log: %v\n", err)
		}
	}
	n, err = lw.bufWriter.Write(p)
	if err != nil {
		return n, fmt.Errorf("log writer failed to write to buffer: %v", err)
	}
	// 优化刷新逻辑
	shouldFlush := len(p) > lw.bufSize/2 ||
		strings.HasSuffix(string(p), "\n") ||
		time.Since(lw.lastFlush) > lw.flushInterval

	if shouldFlush {
		if err = lw.flushBuffer(); err != nil {
			return n, fmt.Errorf("log writer flush failed: %w", err)
		}
		lw.lastFlush = time.Now()
	}
	return n, err
}

// 不加锁的，供内部使用
func (lw *LogWriter) flushBuffer() error {
	if lw.bufWriter == nil {
		return nil
	}
	if err := lw.bufWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	// 同时刷新底层文件
	if lw.file != nil {
		if err := lw.file.Sync(); err != nil {
			return fmt.Errorf("failed to sync file: %w", err)
		}
	}

	return nil
}

// Flush 加锁的，仅供外部使用
func (lw *LogWriter) Flush() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	return lw.flushBuffer()
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	if lw.bufSize > 0 {
		n, err = lw.writeBuffer(p)
		if err == nil {
			return n, nil
		}
		// buffer 方式写入错误，就尝试直接写
		if _, err = lw.write([]byte("write log buffer failed: " + err.Error())); err != nil {
			Console("failed to write log: " + err.Error())
		}
	}
	return lw.write(p)
}

func (lw *LogWriter) openFile() error {
	if lw.bufWriter != nil {
		_ = lw.bufWriter.Flush()
		lw.bufWriter = nil
	}
	if lw.file != nil {
		_ = lw.file.Close()
		lw.file = nil
	}

	newLogFile := lw.currentPath()
	f, err := os.OpenFile(newLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, lw.perm)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	_ = os.Remove(lw.symlink)
	symdir := filepath.Dir(lw.symlink)
	if filepath.Dir(newLogFile) == symdir {
		// 使用相对路径，对于容器挂载路径方式更友好
		if err = os.Chdir(symdir); err == nil {
			_ = os.Symlink(filepath.Base(newLogFile), filepath.Base(lw.symlink))
		}
	} else {
		// 对于虚拟机共享文件夹，可能创建软链会失败，这无影响
		_ = os.Symlink(newLogFile, lw.symlink)
	}

	lw.logName = newLogFile
	lw.file = f //  Do not close this file descriptor
	if lw.bufSize > 0 {
		lw.bufWriter = bufio.NewWriterSize(f, lw.bufSize)
	}
	return nil
}

func (lw *LogWriter) startCleanupRoutine() {
	lw.cleanupTicker = time.NewTicker(24 * time.Hour)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				Console("recover panic routine clearnup: %v", r)
			}
		}()

		for {
			select {
			case <-lw.cleanupTicker.C:
				lw.cleanOldLogs()
			case <-lw.done:
				return
			}
		}
	}()
}

func (lw *LogWriter) cleanOldLogs() {
	Console(os.Remove(lw.symlink))

	cutoff := time.Now().Add(lw.retention)

	_ = filepath.Walk(lw.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), lw.suffix) {
			return nil
		}

		fileDate, err := time.Parse(lw.dateFormat, strings.TrimSuffix(info.Name(), lw.suffix))
		if err == nil && fileDate.Before(cutoff) {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				Console("failed to remove old log %s: %v", path, err)
			}
		}
		return nil
	})
}

func (lw *LogWriter) Shutdown() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	close(lw.done)

	if lw.cleanupTicker != nil {
		lw.cleanupTicker.Stop()
	}
	var errs []error
	if lw.bufWriter != nil {
		if err := lw.bufWriter.Flush(); err != nil {
			errs = append(errs, fmt.Errorf("failed to flush buffer: %w", err))
		}
	}

	if lw.file != nil {
		if err := lw.file.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close file: %w", err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("close log writer errros: %v", errs)
	}
	return nil
}

func redirectStderr(f *os.File) error {
	if err := utils.Dup2(f.Fd(), 2); err != nil {
		return fmt.Errorf("failed to dup fd: %w", err)
	}
	return nil
}
