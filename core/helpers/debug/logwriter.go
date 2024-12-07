package debug

import (
	"bufio"
	"fmt"
	"github.com/aarioai/airis/pkg/arrmap"
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
	logRetentionDays  = 90 // 保留日志的天数
	logFileDateFormat = "2006-01-02"
	logSuffix         = ".bak.log"
)

type LogWriter struct {
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
	lastFlush     time.Time     // 添加最后刷新时间
	flushInterval time.Duration // 添加刷新间隔
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

	// 格式化消息
	msg := arrmap.SprintfArgs(args...)
	if msg == "" {
		return
	}

	// 移除尾部换行符
	if strings.HasSuffix(msg, "\n") {
		msg = msg[:len(msg)-1]
	}
	// 方便运行程序时直接显示
	fmt.Println(msg)
	// 记录进日志，方便通过消息队列通知
	log.Println(msg)
}

// RedirectLog 重定向标准日志和panic输出
func RedirectLog(dir string, perm os.FileMode, bufSize int, symlink ...string) error {
	dir, _ = filepath.Abs(dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %s: %w", dir, err)
	}
	panicLog := path.Join(dir, "panic.log")
	panicFile, err := os.OpenFile(panicLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, perm)
	if err != nil {
		return fmt.Errorf("failed to open crash file: %w", err)
	}
	defer panicFile.Close()

	// 重定向标准错误输出
	if err := redirectStderr(panicFile); err != nil {
		return fmt.Errorf("failed to redirect stderr: %w", err)
	}
	lw, err := NewLogWriter(dir, perm, bufSize, symlink...)
	if err != nil {
		return fmt.Errorf("failed to create log writer: %w", err)
	}
	// 设置日志输出
	log.SetOutput(lw)
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate | log.Lmicroseconds)
	return nil
}

// NewLogWriter 单例模式创建新的日志写入器
// @warn docker 内使用，需要注意单独添加软链接映射，如 -v /var/log/symlink.log:/var/log/symlink.log
func NewLogWriter(dir string, perm os.FileMode, bufSize int, symlinks ...string) (*LogWriter, error) {
	dir, err := ios.MkdirAll(dir, perm)
	if err != nil {
		return nil, err
	}

	symlink := path.Join(dir, "app.log")
	if len(symlinks) > 0 {
		symlink, _, err = ios.PrepareFile(symlinks[0], os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	onceLogWriter.Do(func() {
		logWriter = &LogWriter{
			perm:    perm,
			dir:     dir,
			symlink: symlink,
			bufSize: bufSize,
			done:    make(chan struct{}),
		}
		logWriter.initialize()
		logWriter.startCleanupRoutine()
	})

	// 检查并更新配置
	if needsUpdate := logWriter.updateConfig(dir, symlink, bufSize); needsUpdate {
		logWriter.initialize()
	}

	return logWriter, nil
}
func (lw *LogWriter) initialize() {
	if err := lw.openFile(); err != nil {
		Console("initialize log file: %v", err)
	}

	if lw.bufSize > 0 {
		Console("initialize log writer with buffer size %dB", lw.bufSize)
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
	return path.Join(lw.dir, time.Now().Format(logFileDateFormat)+logSuffix)
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
		if err = lw.Flush(); err != nil {
			return n, fmt.Errorf("log writer flush failed: %w", err)
		}
		lw.lastFlush = time.Now()
	}
	return n, err
}
func (lw *LogWriter) Flush() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	if lw.bufWriter == nil {
		return nil
	}
	return lw.bufWriter.Flush()
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
	// 关闭现有文件和缓冲区
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

	// 对于虚拟机共享文件夹，可能创建软链会失败，这无影响
	_ = os.Remove(lw.symlink)
	Console(os.Symlink(newLogFile, lw.symlink))

	lw.logName = newLogFile
	lw.file = f // 不要关闭 f
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

	cutoff := time.Now().Add(-time.Hour * 24 * time.Duration(logRetentionDays))

	_ = filepath.Walk(lw.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), logSuffix) {
			return nil
		}

		fileDate, err := time.Parse(logFileDateFormat, strings.TrimSuffix(info.Name(), logSuffix))
		if err == nil && fileDate.Before(cutoff) {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				Console("failed to remove old log %s: %v", path, err)
			}
		}
		return nil
	})
}

// Shutdown 正常该日志函数应该是伴随程序终生的，不应该关闭
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

// redirectStderr 重定向标准错误输出
func redirectStderr(f *os.File) error {
	// 使用平台特定的系统调用
	if err := utils.Dup2(f.Fd(), 2); err != nil {
		return fmt.Errorf("failed to dup fd: %w", err)
	}
	return nil
}