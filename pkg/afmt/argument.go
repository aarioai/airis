package afmt

import (
	"fmt"
	"github.com/aarioai/airis/pkg/types"
)

func DefaultIfZero[T types.BasicType](v, defaultV T) T {
	var zero T
	if v == zero {
		return defaultV
	}
	return v
}

// First returns the first argument, or the zero value when the slice is empty. Typically used to read variadic arguments.
func First[T any](args []T) T {
	if len(args) == 0 {
		var zero T
		return zero
	}
	return args[0]
}

// SprintfArgs formats a string from optional arguments.
func SprintfArgs[T any](args []T) string {
	if len(args) == 0 {
		return ""
	}

	format, ok := any(args[0]).(string)
	if !ok {
		format = fmt.Sprint(args[0])
	}

	switch len(args) {
	case 1:
		return format
	case 2:
		return fmt.Sprintf(format, any(args[1]))
	default:
		// Convert the remaining arguments to []any.
		rest := make([]any, len(args)-1)
		for i, arg := range args[1:] {
			rest[i] = any(arg)
		}
		return fmt.Sprintf(format, rest...)
	}
}
