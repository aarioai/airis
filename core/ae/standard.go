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

	ErrEmptyInput       = errors.New("input is empty") // 必需的参数未传递，或者传递为空
	ErrInputWrongLength = errors.New("input wrong length")
	ErrInvalidInput     = errors.New("invalid input") // 存在无效的参数

)

// *Error
// ErrorXXX/ErrXXX  都应被视为常量，不应修改
var (
	ErrorEmptyInput       = NewError(ErrEmptyInput).WithCaller(2).Lock()
	ErrorInvalidInput     = NewError(ErrInvalidInput).WithCaller(2).Lock()
	ErrorInputWrongLength = NewError(ErrInputWrongLength).WithCaller(2).Lock()
	ErrorDeprecated       = NewError(ErrDeprecated).WithCaller(2).Lock()
)

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
func ErrInvalid(value any, name ...string) error {
	return errors.New("invalid input  " + packArg(value, name))
}

func ErrorInvalid(value any, name ...string) *Error {
	return NewE("invalid input " + packArg(value, name))
}
