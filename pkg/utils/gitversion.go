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

// GitVersion returns the Git commit hash.
func GitVersion() string {
	once.Do(func() {
		gitHash = loadGitHash()
	})
	return gitHash
}

// loadGitHash loads the Git hash appended to the executable.
// Expected format: a newline plus a 40-byte hash (left-padded with spaces; TrimLeft after read).
// printf '\n%s' "$(git rev-parse HEAD 2>/dev/null)" >> main
// Even on read errors the hash can be verified against real git, so impact is minimal.
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

// lookExecPath returns the real path of the current executable.
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

// isValidGitHash reports whether hash is a valid Git hash.
func isValidGitHash(hash []byte) bool {
	for _, h := range hash {
		if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'z')) {
			return false
		}
	}
	return true
}
