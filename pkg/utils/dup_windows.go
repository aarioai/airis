//go:build windows
// +build windows

package utils

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func Dup2(oldfd uintptr, newfd int) error {
	// 参数验证增强
	if oldfd == uintptr(windows.InvalidHandle) || newfd < 0 {
		return fmt.Errorf("invalid file descriptor: oldfd=%d, newfd=%d", oldfd, newfd)
	}

	// 如果新旧句柄相同，直接返回
	if uintptr(newfd) == oldfd {
		return nil
	}

	// 获取当前进程句柄
	currentProcess := windows.CurrentProcess()

	// 复制句柄
	var newHandle windows.Handle
	err := windows.DuplicateHandle(
		currentProcess,
		windows.Handle(oldfd),
		currentProcess,
		&newHandle,
		0,
		false,
		windows.DUPLICATE_SAME_ACCESS,
	)
	if err != nil {
		return fmt.Errorf("failed to duplicate handle: %w", err)
	}

	// 确保在发生错误时关闭新句柄
	defer func() {
		if err != nil {
			_ = windows.CloseHandle(newHandle)
		}
	}()

	// 关闭旧的句柄（如果存在且有效）
	oldHandle := windows.Handle(newfd)
	if oldHandle != windows.InvalidHandle {
		if err := windows.CloseHandle(oldHandle); err != nil {
			return fmt.Errorf("failed to close existing handle: %w", err)
		}
	}

	// 设置新的句柄值
	*(*windows.Handle)(unsafe.Pointer(uintptr(newfd))) = newHandle

	return nil
}
