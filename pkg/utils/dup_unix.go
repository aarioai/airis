//go:build !windows
// +build !windows

package utils

import (
	"fmt"
	"syscall"
)

func Dup2(oldfd uintptr, newfd int) error {
	if oldfd == 0 || newfd < 0 {
		return fmt.Errorf("invalid file descriptor: oldfd=%d, newfd=%d", oldfd, newfd)
	}

	// If the old and new file descriptors are the same, return immediately.
	if uintptr(newfd) == oldfd {
		return nil
	}

	// If dup3 is unavailable or fails, fall back to dup2.
	return syscall.Dup2(int(oldfd), newfd)
}
