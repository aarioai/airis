package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	defaultLogBufSize = 2048 // 默认缓冲区 2KB
	logRetentionDays  = 90   // 保留日志的天数
	logFileDateFormat = "2006-01-02"
	logSuffix         = ".bak.log"
)

type LogWriter struct {
	mode          os.FileMode
	filename      string
	file          *os.File
	logName       string
	mu            sync.RWMutex
	cleanupTicker *time.Ticker
	bufSize       int
	bufWriter     *bufio.Writer
	done          chan struct{}
}

var (
	onceLogWriter sync.Once
	logWriter     *LogWriter
)

// RedirectLog 重定向标准日志和panic输出
func RedirectLog(logfile, crashfile string, mode, dirmode os.FileMode) error {
	for _, dir := range []string{path.Dir(logfile), path.Dir(crashfile)} {
		if err := os.MkdirAll(dir, dirmode); err != nil {
			return fmt.Errorf("failed to create log directory: %s: %w", dir, err)
		}
	}

	absLogfile, err := filepath.Abs(logfile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	crashFile, err := os.OpenFile(crashfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open crash file: %w", err)
	}
	defer crashFile.Close()

	// 重定向标准错误输出
	if err := redirectStderr(crashFile); err != nil {
		return fmt.Errorf("failed to redirect stderr: %w", err)
	}

	// 设置日志输出
	log.SetOutput(NewLogWriter(absLogfile, mode))
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate | log.Lmicroseconds)
	return nil
}

// NewLogWriter 单例模式创建新的日志写入器
func NewLogWriter(logfile string, mod os.FileMode) *LogWriter {
	onceLogWriter.Do(func() {
		logWriter = &LogWriter{
			filename: logfile,
			bufSize:  defaultLogBufSize,
			done:     make(chan struct{}),
		}
		if err := logWriter.rotate(); err != nil {
			log.Printf("initialize log file: %v", err)
		}
		logWriter.startCleanupRoutine()
	})
	return logWriter
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	lw.mu.RLock()
	currentFile := lw.logName
	lw.mu.RUnlock()

	newLogFile := path.Join(path.Dir(lw.filename), time.Now().Format(logFileDateFormat)+logSuffix)
	if newLogFile != currentFile {
		lw.mu.Lock()
		if newLogFile != lw.logName {
			if err := lw.rotate(); err != nil {
				log.Printf("failed to rotate log: %v\n", err)
			}
		}
		lw.mu.Unlock()
	}

	lw.mu.RLock()
	defer lw.mu.RUnlock()

	if lw.bufWriter == nil {
		return 0, fmt.Errorf("log writer not initialized")
	}
	n, err = lw.bufWriter.Write(p)
	if err == nil && (len(p) > lw.bufSize/2 || strings.HasSuffix(string(p), "\n")) {
		err = lw.bufWriter.Flush()
	}
	return
}

func (lw *LogWriter) rotate() error {
	newLogFile := path.Join(path.Dir(lw.filename), time.Now().Format(logFileDateFormat)+logSuffix)
	f, err := os.OpenFile(newLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, lw.mode)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}
	if lw.bufWriter != nil {
		_ = lw.bufWriter.Flush()
	}
	if lw.file != nil {
		_ = lw.file.Close()
	}
	_ = os.Remove(lw.filename)
	if err := os.Symlink(newLogFile, lw.filename); err != nil {
		return fmt.Errorf("failed to create new symlink: %w", err)
	}

	lw.logName = newLogFile
	lw.file = f
	lw.bufWriter = bufio.NewWriterSize(f, lw.bufSize)
	return nil
}

func (lw *LogWriter) startCleanupRoutine() {
	lw.cleanupTicker = time.NewTicker(24 * time.Hour)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recover panic routine clearnup: %v", r)
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
	dir := path.Dir(lw.filename)
	cutoff := time.Now().Add(-time.Hour * 24 * time.Duration(logRetentionDays))

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), logSuffix) {
			return nil
		}

		fileDate, err := time.Parse(logFileDateFormat, strings.TrimSuffix(info.Name(), logSuffix))
		if err == nil && fileDate.Before(cutoff) {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				log.Printf("failed to remove old log %s: %v", path, err)
			}
		}
		return nil
	})
}

func (lw *LogWriter) Close() error {
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
	if err := Dup2(f.Fd(), 2); err != nil {
		return fmt.Errorf("failed to dup fd: %w", err)
	}
	return nil
}
