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

// First 获取第一个参数，如果切片为空则返回零值。一般用于动态参数获取
func First[T any](args []T) T {
	if len(args) == 0 {
		var zero T
		return zero
	}
	return args[0]
}

// SprintfArgs 针对可选参数，格式化字符串
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
		// 将剩余参数转换为 []any 类型
		rest := make([]any, len(args)-1)
		for i, arg := range args[1:] {
			rest[i] = any(arg)
		}
		return fmt.Sprintf(format, rest...)
	}
}
