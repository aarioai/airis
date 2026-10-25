package afmt

import (
	"fmt"
	"reflect"
)

// ErrmsgSideEffect is the error text for a function that mutates an input it should not change (those are usually marked @note slice arguments are safe).
// Typically used as an error message in a function or a unit test.
func ErrmsgSideEffect(variable any) string {
	return fmt.Sprintf("improper use of %s occurs unintended side effects", reflect.TypeOf(variable))
}
