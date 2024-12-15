package ae

import (
	"github.com/aarioai/airis/pkg/afmt"
)

func First(es ...*Error) *Error {
	return afmt.First(es)
}
func FirstErr(errs ...error) error {
	return afmt.First(errs)
}

// PanicOn 如果存在服务器错误则触发 panic
func PanicOn(es ...*Error) {
	if e := First(es...); e != nil {
		panic(e.Text())
	}
}

// PanicOnErrors 断言检查标准错误，如果存在错误则触发 panic
func PanicOnErrs(errs ...error) {
	if e := FirstErr(errs...); e != nil {
		panic(e.Error())
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
