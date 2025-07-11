package ae

import (
	"errors"
	"fmt"
)

// 定义标准错误，节省内存开销以及方便通用处理
// 这里的错误，都是代码内部错误，不对外暴露

// 标准错误
// ErrorXXX/ErrXXX  都应被视为常量，不应修改
var (
	ErrNotImplemented     = errors.New("not implemented") // 未实现的功能
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrDeprecated         = errors.New("deprecated") // 已废弃的功能

	ErrInvalidInput     = errors.New("invalid input") // 存在无效的参数
	ErrInputWrongLength = errors.New("input wrong length")
	ErrEmptyInput       = errors.New("input is empty") // 必需的参数未传递，或者传递为空
	ErrInputTooLong     = errors.New("input is too long")
	ErrInputTooShort    = errors.New("input is too short")
	ErrInputTooBig      = errors.New("input is too big")
	ErrInputTooSmall    = errors.New("input is too small")

	ErrNotFound = errors.New("not found")
)

// *Error
// ErrorXXX/ErrXXX  都应被视为常量，不应修改
var (
	ErrorDeprecated = NewErr(ErrDeprecated).WithCaller(2).Lock()

	ErrorInvalidInput     = NewErr(ErrInvalidInput).WithCaller(2).Lock()
	ErrorInputWrongLength = NewErr(ErrInputWrongLength).WithCaller(2).Lock()
	ErrorEmptyInput       = NewErr(ErrEmptyInput).WithCaller(2).Lock()
	ErrorInputTooLong     = NewErr(ErrInputTooLong).WithCaller(2).Lock()
	ErrorInputTooShort    = NewErr(ErrInputTooShort).WithCaller(2).Lock()
	ErrorInputTooBig      = NewErr(ErrInputTooBig).WithCaller(2).Lock()
	ErrorInputTooSmall    = NewErr(ErrInputTooSmall).WithCaller(2).Lock()
)

func ErrInvalid(value any, name ...string) error {
	return errors.New("invalid input  " + packArg(value, name))
}

func ErrorInvalid(value any, name ...string) *Error {
	return NewErrorf("invalid input %s", packArg(value, name))
}
func packArg(value any, name []string) string {
	var s string
	if len(name) > 0 {
		s = name[0]
	}
	if value != nil {
		s += fmt.Sprintf(": %v", value)
	}
	return s
}
