package utils

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

const (
	gitHashLength = 40
	timeFormat    = "2006-01-02 15:04:05"
)

var (
	startingSteps int32
	gitHash       string
	once          sync.Once
)

// GitVersion 获取Git版本哈希值
func GitVersion() string {
	once.Do(func() {
		gitHash = loadGitHash()
	})
	return gitHash
}

// loadGitHash 从可执行文件中加载Git哈希值
func loadGitHash() string {
	execPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}

	file, err := os.OpenFile(execPath, os.O_RDONLY, 0666)
	if err != nil {
		return ""
	}
	defer file.Close()

	finfo, err := file.Stat()
	if err != nil {
		return ""
	}

	hash := make([]byte, gitHashLength)
	if _, err := file.ReadAt(hash, finfo.Size()-gitHashLength-1); err != nil {
		return ""
	}

	if !isValidGitHash(hash) {
		return ""
	}

	return string(hash)
}

// isValidGitHash 验证是否为有效的Git哈希值
func isValidGitHash(hash []byte) bool {
	for _, h := range hash {
		if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'z')) {
			return false
		}
	}
	return true
}

// TraceStartup 记录启动步骤
func TraceStartup(msg ...string) {
	id := atomic.AddInt32(&startingSteps, 1)
	now := formatCurrentTime()

	message := fmt.Sprintf("%s starting %d", now, id)
	if len(msg) > 0 {
		message += " " + msg[0]
	}

	fmt.Println(message)
}

// formatCurrentTime 格式化当前时间，包含微秒
func formatCurrentTime() string {
	n := time.Now()
	return fmt.Sprintf("%s.%d",
		n.Format(timeFormat),
		n.UnixMicro(),
	)
}
