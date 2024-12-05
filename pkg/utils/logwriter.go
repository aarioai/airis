package utils

import (
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
	defaultFileMode  = 0666
	logRetentionDays = 90
	dateFormat       = "2006-01-02"
	logSuffix        = ".bak.log"
)

var (
	onceLogWriter sync.Once
	logWriter     *LogWriter
)

type LogWriter struct {
	filename      string
	file          *os.File
	logName       string
	mu            sync.RWMutex
	cleanupTicker *time.Ticker
}

func RedirectLog(logfile, crashfile string, mode os.FileMode) error {
	if err := os.MkdirAll(path.Dir(logfile), mode); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	if err := os.MkdirAll(path.Dir(crashfile), mode); err != nil {
		return fmt.Errorf("failed to create crash directory: %w", err)
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
    log.SetOutput(NewLogWriter(absLogfile))
    log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate | log.Lmicroseconds)
	return nil
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	dir := path.Dir(lw.filename)
	newLogFile := path.Join(dir, time.Now().Format(dateFormat)+logSuffix)

	lw.mu.RLock()
	needRotate := newLogFile != lw.logName
	lw.mu.RUnlock()

	if needRotate {
		lw.mu.Lock()
		if newLogFile != lw.logName {
			if err := lw.rotateLog(newLogFile); err != nil {
				log.Printf("Failed to rotate log: %v\n", err)
			}
		}
		lw.mu.Unlock()
	}

	lw.mu.RLock()
	defer lw.mu.RUnlock()

	if lw.file != nil {
		return lw.file.Write(p)
	}
	return 0, fmt.Errorf("log file not initialized")
}

func (lw *LogWriter) rotateLog(newLogFile string) error {
	f, err := os.OpenFile(newLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, defaultFileMode)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	if _, err := f.WriteString("\n\n\n"); err != nil {
		f.Close()
		return fmt.Errorf("failed to write header: %w", err)
	}

	if lw.file != nil {
		if err := lw.file.Close(); err != nil {
			log.Printf("Failed to close old log file: %v", err)
		}
	}

	lw.logName = newLogFile
	lw.file = f

	if err := os.Remove(lw.filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove old symlink: %w", err)
	}

	return os.Symlink(newLogFile, lw.filename)
}

func (lw *LogWriter) startCleanupRoutine() {
	lw.cleanupTicker = time.NewTicker(24 * time.Hour)
	go func() {
		for range lw.cleanupTicker.C {
			lw.cleanOldLogs()
		}
	}()
}

func (lw *LogWriter) cleanOldLogs() {
	dir := path.Dir(lw.filename)
	now := time.Now()

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to read log directory: %v", err)
		return
	}

	cutoff := now.Add(-time.Hour * 24 * time.Duration(logRetentionDays))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), logSuffix) {
			fileDate, err := time.Parse(dateFormat, strings.TrimSuffix(entry.Name(), logSuffix))
			if err != nil {
				continue
			}

			if fileDate.Before(cutoff) {
				oldLogPath := path.Join(dir, entry.Name())
				if err := os.Remove(oldLogPath); err != nil && !os.IsNotExist(err) {
					log.Printf("Failed to remove old log %s: %v", oldLogPath, err)
				}
			}
		}
	}
}

func NewLogWriter(logfile string) *LogWriter {
	onceLogWriter.Do(func() {
		logWriter = &LogWriter{
			filename: logfile,
		}
		logWriter.startCleanupRoutine()
	})
	return logWriter
}

func (lw *LogWriter) Close() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.cleanupTicker != nil {
		lw.cleanupTicker.Stop()
	}

	if lw.file != nil {
		return lw.file.Close()
	}
	return nil
}


// redirectStderr 重定向标准错误输出
func redirectStderr(f *os.File) error {
    // 使用平台特定的系统调用
    if err := dupFd(f.Fd(), 2); err != nil {
        return fmt.Errorf("failed to dup fd: %w", err)
    }
    return nil
}

// dupFd 复制文件描述符
// 根据不同操作系统使用不同的实现
func dupFd(fd uintptr, newfd int) error {
    return dup2(fd, newfd)
}