package atype

// Deref 是一个通用的指针解引用函数，如果指针为 nil 则返回类型的零值
// T 必须是基本数值类型: uint64, uint32, uint, uint16, uint8, int64, int32, int, int16, int8, float64, float32, string
// @example a:=100; var b *int; Deref(&a)  Deref(b)
func Deref[T uint64 | uint32 | uint | uint16 | uint8 | int64 | int32 | int | int16 | int8 | float64 | float32 | string](n *T) T {
	if n == nil {
		var zero T
		return zero
	}
	return *n
}

// DerefInt24 专门处理 Int24 类型的指针解引用
func DerefInt24(n *Uint24) Uint24 {
	if n == nil {
		return 0
	}
	return *n
}

// DerefUint24 专门处理 Uint24 类型的指针解引用
func DerefUint24(n *Uint24) Uint24 {
	if n == nil {
		return 0
	}
	return *n
}
