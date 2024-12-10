package conv

import (
	"github.com/aarioai/airis/core/atype"
	"strconv"
)

// ParseInt 通用数字解析函数
// @example ParseInt[int8]("123", 10, 8)    ParseInt[uint8]("123", 10, 8)
func ParseInt[T ~int8 | ~int16 | ~int32 | ~int64](s string, base int, bitSize int) T {
	v, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		var zero T
		return zero
	}
	return T(v)
}

// ParseUint 通用无符号数字解析函数
// @example ParseUint[uint8]("123", 10, 8)    ParseUint[uint16]("123", 10, 8)
func ParseUint[T ~uint8 | ~uint16 | ~uint32 | ~uint64](s string, base int, bitSize int) T {
	u, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		var zero T
		return zero
	}
	return T(u)
}
 


// IsUnsigned 判断类型是否为无符号整数
func IsUnsigned[T any]() bool {
	var zero T
	switch any(zero).(type) {
	case uint8, uint16, atype.Uint24, uint32, uint64:
		return true
	default:
		return false
	}
}