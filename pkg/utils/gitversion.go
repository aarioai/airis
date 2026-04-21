package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

const (
	gitHashLength = 40
)

var (
	gitHash string
	once    sync.Once
)

func GitVersion() string {
	once.Do(func() {
		gitHash = loadGitHash()
	})
	return gitHash
}

// loadGitHash
// 预期格式: 40字节的哈希值（不足部分用空格填充左边，实际读取后TrimLeft）
// printf '\n%s' "$(git rev-parse HEAD 2>/dev/null)" >> main
// 即使获取错误，可以对实际git进行核对，因此几乎没有任何影响
func loadGitHash() string {
	var (
		execPath string
		file     *os.File
		finfo    os.FileInfo
		err      error
	)

	if execPath, err = lookExecPath(); err != nil {
		return ""
	}

	if file, err = os.Open(execPath); err != nil {
		return ""
	}
	defer file.Close()

	if finfo, err = file.Stat(); err != nil {
		return ""
	}
	size := finfo.Size()
	if size <= gitHashLength {
		return ""
	}
	hash := make([]byte, gitHashLength)
	if _, err = file.ReadAt(hash, size-gitHashLength-1); err != nil {
		return ""
	}
	if hash[0] != '\n' {
		return ""
	}
	hash = hash[1:]
	hash = bytes.TrimLeft(hash, " ")
	if !isValidGitHash(hash) {
		return ""
	}
	return string(hash)
}

// lookExecPath returns the real path of the current executable
func lookExecPath() (string, error) {
	// Prefer os.Executable() as it's more reliable
	execPath, err := os.Executable()
	if err != nil {
		// Fallback to exec.LookPath
		execPath, err = exec.LookPath(os.Args[0])
		if err != nil {
			return "", fmt.Errorf("failed to get executable path: %w", err)
		}
	}

	// Resolve symlinks if any
	var realPath string
	if realPath, err = filepath.EvalSymlinks(execPath); err == nil {
		execPath = realPath
	}

	return execPath, nil
}

func isValidGitHash(hash []byte) bool {
	for _, h := range hash {
		if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'z')) {
			return false
		}
	}
	return true
}
