package atype

import "github.com/aarioai/airis/pkg/types"

func ToInt24(s string) Int24 {
	v, _ := ParseInt24(s)
	return v
}
func ToUint24(s string) Uint24 {
	v, _ := ParseUint24(s)
	return v
}
func ParseInt24(s string) (Int24, error) {
	v, err := types.ParseSigned[int32](s, 10, 24)
	return Int24(v), err
}
func ParseUint24(s string) (Uint24, error) {
	v, err := types.ParseUnsigned[uint32](s, 10, 24)
	return Uint24(v), err
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
