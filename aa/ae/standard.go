package ae

import (
	"errors"
	"fmt"
)

// Standard errors, defined once to save memory and share handling.
// These are internal errors and are not exposed outside the code.

// Standard errors.
// ErrorXXX/ErrXXX must be treated as constants and must not be modified.
var (
	ErrNotImplemented     = errors.New("not implemented") // feature not implemented
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrDeprecated         = errors.New("deprecated") // deprecated feature

	ErrInvalidInput     = errors.New("invalid input") // an invalid parameter is present
	ErrInputWrongLength = errors.New("input wrong length")
	ErrEmptyInput       = errors.New("input is empty") // a required parameter was omitted or passed empty
	ErrInputTooLong     = errors.New("input is too long")
	ErrInputTooShort    = errors.New("input is too short")
	ErrInputTooBig      = errors.New("input is too big")
	ErrInputTooSmall    = errors.New("input is too small")

	ErrNotFound = errors.New("not found")
)

// *Error
// ErrorXXX/ErrXXX must be treated as constants and must not be modified.
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
