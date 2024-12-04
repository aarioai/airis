package conv

import (
	"strconv"

	"github.com/aarioai/airis/core/atype"
)

// Atoi 通用数字解析函数
// @example:
// Atoi[int8]("123",10,8)
// Atoi[uint8]("123",10,8)
func Atoi[T ~int8 | ~int16 | ~int32 | ~uint8 | ~uint16 | ~uint32](s string, base int, bitSize int) (T, error) {
	var v int64
	var err error

	if IsUnsigned[T]() {
		u, err := strconv.ParseUint(s, base, bitSize)
		if err != nil {
			return 0, err
		}
		v = int64(u)
	} else {
		v, err = strconv.ParseInt(s, base, bitSize)
		if err != nil {
			return 0, err
		}
	}

	return T(v), nil
}

// IsUnsigned 判断类型是否为无符号整数
func IsUnsigned[T any]() bool {
	var zero T
	switch any(zero).(type) {
	case uint8, uint16, uint32, uint64, atype.Uint24:
		return true
	default:
		return false
	}
}
