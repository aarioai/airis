package atype

import (
	"golang.org/x/exp/constraints"
	"math"
	"strconv"
)

const maxSmallNumber = math.MaxUint8

// 为了提高性能，可以添加常用数字的字符串缓存。小数字在query string 或body传参中会很常见
var (
	smallNumbers [maxSmallNumber + 1]string // 暂时保持 uint8 最高范围
)

func init() {
	for i := 0; i < maxSmallNumber; i++ {
		smallNumbers[i] = strconv.Itoa(i)
	}
}

// FormatUint8 直接查表
// 小数字在query string 或body传参中会很常见
func FormatUint8(v uint8) string {
	return smallNumbers[v]
}

// FormatInt 将int64转换为字符串，小数字直接使用查表法
func FormatInt[T constraints.Signed](value T) string {
	v := int64(value)
	// 小数字使用查表法
	if v >= 0 && v <= maxSmallNumber {
		return smallNumbers[v]
	}
	return strconv.FormatInt(int64(v), 10)
}

// FormatUint 将uint64转换为字符串，小数字直接使用查表法
func FormatUint[T constraints.Unsigned](value T) string {
	v := uint64(value)
	// 小数字使用查表法
	if v <= maxSmallNumber {
		return smallNumbers[v]
	}
	return strconv.FormatUint(v, 10)
}

// FormatFloat 将float转换为字符串，整数部分使用FormatInt，小数部分使用strconv.FormatFloat
func FormatFloat[T constraints.Float](value T, bitSize int) string {
	v := float64(value)
	intPart := int64(v)
	// 处理int性能更好
	if v == float64(intPart) {
		return FormatInt(intPart)
	}
	return strconv.FormatFloat(v, 'f', -1, bitSize)
}
func ParseBool(s string) bool {
	v, _ := strconv.ParseBool(s)
	return v
}
func ParseInt(s string) int {
	return parseInt[int](s, 10, 32)
}
func ParseInt8(s string) int8 {
	return parseInt[int8](s, 10, 8)
}
func ParseInt16(s string) int16 {
	return parseInt[int16](s, 10, 16)
}
func ParseInt32(s string) int32 {
	return parseInt[int32](s, 10, 32)
}
func ParseInt64(s string) int64 {
	return parseInt[int64](s, 10, 64)
}

func ParseUint(s string) uint {
	return parseUint[uint](s, 10, 32)
}
func ParseUint8(s string) uint8 {
	return parseUint[uint8](s, 10, 8)
}
func ParseUint16(s string) uint16 {
	return parseUint[uint16](s, 10, 16)
}
func ParseUint32(s string) uint32 {
	return parseUint[uint32](s, 10, 32)
}
func ParseUint64(s string) uint64 {
	return parseUint[uint64](s, 10, 64)
}
func ParseFloat32(s string) float32 {
	result, _ := strconv.ParseFloat(s, 32)
	return float32(result)
}
func ParseFloat64(s string) float64 {
	result, _ := strconv.ParseFloat(s, 64)
	return result
}

// parseInt 通用数字解析函数
// @example parseInt[int8]("123", 10, 8)    parseInt[uint8]("123", 10, 8)
func parseInt[T ~int8 | ~int16 | ~int32 | ~int | ~int64](s string, base int, bitSize int) T {
	v, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		var zero T
		return zero
	}
	return T(v)
}

// parseUint 通用无符号数字解析函数
// @example parseUint[uint8]("123", 10, 8)    parseUint[uint16]("123", 10, 8)
func parseUint[T ~uint8 | ~uint16 | ~uint32 | ~uint | ~uint64](s string, base int, bitSize int) T {
	u, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		var zero T
		return zero
	}
	return T(u)
}

// IsUnsigned 判断类型是否为无符号整数
func IsUnsigned[T any](_ T) bool {
	var zero T
	switch any(zero).(type) {
	case uint8, uint16, Uint24, uint32, uint64:
		return true
	default:
		return false
	}
}
