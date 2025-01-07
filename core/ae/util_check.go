package ae

import (
	"fmt"
	"github.com/aarioai/airis/pkg/afmt"
	"time"
)

func First(es ...*Error) *Error {
	return afmt.First(es)
}
func FirstErr(errs ...error) error {
	return afmt.First(errs)
}

func Panic(format string, args ...any) {
	now := time.Now().Format("2006-01-02 15:04:05")
	msg := now + " [panic] " + fmt.Sprintf(format, args...)
	fmt.Println(msg) // convenient for docker/local debugging
	panic(msg)
}

// PanicOn 如果存在服务器错误则触发 panic
func PanicOn(es ...*Error) {
	if e := First(es...); e != nil {
		Panic(e.Text())
	}
}

// PanicOnErrors 断言检查标准错误，如果存在错误则触发 panic
func PanicOnErrs(errs ...error) {
	if e := FirstErr(errs...); e != nil {
		Panic(e.Error())
	}
}
func StopOnFirstError(callables ...func() *Error) *Error {
	for _, callable := range callables {
		if err := callable(); err != nil {
			return err
		}
	}
	return nil
}
func StopOnFirstErr(callables ...func() error) error {
	for _, callable := range callables {
		if err := callable(); err != nil {
			return err
		}
	}
	return nil
}
