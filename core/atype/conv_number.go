package atype

import "strconv"

// FormatInt 将int64转换为字符串，小数字直接使用查表法
func FormatInt(v int64) string {
	// 小数字使用查表法；小数字在query string 或body传参中会很常见
	if v >= 0 && v < 100 {
		return getSmallNumberString(int(v))
	}
	return strconv.FormatInt(v, 10)
}

// FormatUint 将uint64转换为字符串，小数字直接使用查表法
func FormatUint(v uint64) string {
	// 小数字使用查表法
	if v < 100 {
		return getSmallNumberString(int(v))
	}
	return strconv.FormatUint(v, 10)
}

// FormatFloat 将float64转换为字符串，整数部分使用FormatInt，小数部分使用strconv.FormatFloat
func FormatFloat(v float64, bitSize int) string {
	// 整数部分处理
	if v == float64(int64(v)) {
		return FormatInt(int64(v))
	}
	// 使用 strconv.FormatFloat
	return strconv.FormatFloat(v, 'f', -1, bitSize)
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
