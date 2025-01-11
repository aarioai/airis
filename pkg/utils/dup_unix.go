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

	// 如果新旧文件描述符相同，直接返回
	if uintptr(newfd) == oldfd {
		return nil
	}

	// 如果 dup3 不可用或失败，回退到 dup2
	return syscall.Dup2(int(oldfd), newfd)
}
