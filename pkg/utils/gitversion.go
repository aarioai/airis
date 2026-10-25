package utils

import (
	"os"
	"os/exec"
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

// loadGitHash loads the Git hash from the executable.
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

// isValidGitHash reports whether hash is a valid Git hash.
func isValidGitHash(hash []byte) bool {
	for _, h := range hash {
		if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'z')) {
			return false
		}
	}
	return true
}
