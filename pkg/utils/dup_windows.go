//go:build windows
// +build windows

package utils

import (
	"fmt"
	"golang.org/x/sys/windows"
)

func Dup2(oldfd uintptr, newfd int) error {
	if oldfd == uintptr(windows.InvalidHandle) || newfd < 0 {
		return fmt.Errorf("invalid file descriptor: oldfd=%d, newfd=%d", oldfd, newfd)
	}
	// Returns on same
	if uintptr(newfd) == oldfd {
		return nil
	}

	currentProcess := windows.CurrentProcess()

	// Map Unix-style file descriptors to Windows standard handles
	var stdHandle uint32
	switch newfd {
	case 0:
		stdHandle = windows.STD_INPUT_HANDLE
	case 1:
		stdHandle = windows.STD_OUTPUT_HANDLE
	case 2:
		stdHandle = windows.STD_ERROR_HANDLE
	default:
		return fmt.Errorf("only standard FDs (0,1,2) are supported on Windows")
	}

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

	// Set the new handle as the standard handle
	if err := windows.SetStdHandle(stdHandle, newHandle); err != nil {
		windows.CloseHandle(newHandle)
		return fmt.Errorf("failed to set std handle: %w", err)
	}

	return nil
}
