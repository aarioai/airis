package ae

import (
	"errors"
	"fmt"
)

// 定义标准错误，节省内存开销以及方便通用处理
// 这里的错误，都是代码内部错误，不对外暴露

// 标准错误
var (
	ErrNotImplemented = errors.New("not implemented") // 未实现的功能

	ErrEmptyInput    = errors.New("input is empty") // 必需的参数未传递，或者传递为空
	ErrInvalidInput  = errors.New("invalid input")  // 存在无效的参数
	ErrInputTooLong  = errors.New("input is too long")
	ErrInputTooShort = errors.New("input is too short")
	ErrDeprecated    = errors.New("deprecated") // 已废弃的功能
)

// *Error
var (
	ErrorEmptyInput    = NewError(ErrEmptyInput).WithCaller(2).Lock()
	ErrorInvalidInput  = NewError(ErrInvalidInput).WithCaller(2).Lock()
	ErrorInputTooLong  = NewError(ErrInputTooLong).WithCaller(2).Lock()
	ErrorInputTooShort = NewError(ErrInputTooShort).WithCaller(2).Lock()
	ErrorDeprecated    = NewError(ErrDeprecated).WithCaller(2).Lock()
)

func ErrInvalid(name string, value ...any) error {
	if len(value) > 0 {
		return fmt.Errorf("invalid input %s: %v", name, value[0])
	}
	return errors.New("invalid input " + name)
}

func ErrorInvalid(name string, value ...any) *Error {
	if len(value) > 0 {
		return NewE(fmt.Sprintf("invalid input %s: %v", name, value[0]))
	}
	return NewE("invalid input " + name)
}
