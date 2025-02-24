package types

import "golang.org/x/exp/constraints"

// IsUnsigned 判断类型是否为无符号整数
func IsUnsigned[T constraints.Integer](_ T) bool {
	var zero T
	switch any(zero).(type) {
	case uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

// Deref 是一个通用的指针解引用函数，如果指针为 nil 则返回类型的零值
// T 必须是基本数值类型: uint64, uint32, uint, uint16, uint8, int64, int32, int, int16, int8, float64, float32, string
// @example a:=100; var b *int; Deref(&a)  Deref(b)
func Deref[T Number | ~string](n *T) T {
	if n == nil {
		var zero T
		return zero
	}
	return *n
}
