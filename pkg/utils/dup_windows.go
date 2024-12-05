//go:build windows
package utils

import "syscall"

func dup2(oldfd uintptr, newfd int) error {
    return syscall.Dup2(syscall.Handle(oldfd), syscall.Handle(newfd))
}