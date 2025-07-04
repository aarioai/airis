package ae

import (
	"errors"
	"fmt"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/utils"
	"strings"
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

func PanicWithCaller(callerSkip int, msg string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	msg = now + " [panic] " + utils.Caller(callerSkip) + " " + msg
	fmt.Println(msg) // convenient for docker/local debugging
	panic(msg)
}

func Panic(msg string) {
	PanicWithCaller(2, msg)
}

func PanicF(format string, args ...any) {
	PanicWithCaller(2, fmt.Sprintf(format, args...))
}

func PanicOn(es ...*Error) {
	if e := First(es...); e != nil {
		PanicWithCaller(2, e.String())
	}
}

func PanicOnErrs(errs ...error) {
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

func Wrap(head string, err error, tails ...string) error {
	if err == nil {
		return nil
	}
	tail := ""
	if len(tails) > 0 {
		tail = " " + strings.Join(tails, " ")
	}
	return errors.New(head + " " + err.Error() + tail)
}
