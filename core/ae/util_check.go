package ae

import (
	"errors"
	"fmt"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/utils"
	"time"
)

func First(es ...*Error) *Error {
	return afmt.First(es)
}
func FirstError(errs ...error) error {
	return afmt.First(errs)
}
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func PanicWithCaller(callerSkip int, format string, args ...any) {
	now := time.Now().Format("2006-01-02 15:04:05")
	msg := now + " [panic] " + utils.Caller(callerSkip) + " " + fmt.Sprintf(format, args...)
	fmt.Println(msg) // convenient for docker/local debugging
	panic(msg)
}

func Panic(format string, args ...any) {
	PanicWithCaller(2, format, args...)
}

// PanicOn 如果存在服务器错误则触发 panic
func PanicOn(es ...*Error) {
	if e := First(es...); e != nil {
		PanicWithCaller(2, e.String())
	}
}

// PanicOnErrors 断言检查标准错误，如果存在错误则触发 panic
func PanicOnErrors(errs ...error) {
	if e := FirstError(errs...); e != nil {
		PanicWithCaller(2, e.Error())
	}
}
func StopOnFirstE(callables ...func() *Error) *Error {
	for _, callable := range callables {
		if err := callable(); err != nil {
			return err
		}
	}
	return nil
}
func StopOnFirstError(callables ...func() error) error {
	for _, callable := range callables {
		if err := callable(); err != nil {
			return err
		}
	}
	return nil
}
