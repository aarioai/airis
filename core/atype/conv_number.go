package atype

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"golang.org/x/exp/constraints"
)

const maxSmallNumber = math.MaxUint8

// 为了提高性能，可以添加常用数字的字符串缓存。小数字在query string 或body传参中会很常见
var (
	smallNumbers    [maxSmallNumber + 1]string // 暂时保持 uint8 最高范围
	smallNumbersMap = make(map[string]int, len(smallNumbers))
)

func init() {
	for i := 0; i <= maxSmallNumber; i++ {
		s := strconv.Itoa(i)
		smallNumbers[i] = s
		smallNumbersMap[s] = i
	}
}

// ConvertBase 数字型字符串进制转换
func ConvertBase(s string, fromBase, toBase int) (string, error) {
	if s == "" {
		return "", errors.New("string is empty")
	}
	if s == "0" {
		return "0", nil
	}
	if s[0] == '-' {
		var v int64
		var err error
		if fromBase == 10 {
			v, err = ParseInt64(s)
		} else {
			v, err = strconv.ParseInt(s, fromBase, 64)
		}

		if err != nil {
			return "", err
		}
		if toBase == 10 {
			return FormatInt(v), nil
		}
		return strconv.FormatInt(v, toBase), nil
	}
	// unsigned
	var v uint64
	var err error
	if fromBase == 10 {
		v, err = ParseUint64(s)
	} else {
		v, err = strconv.ParseUint(s, fromBase, 64)
	}
	if err != nil {
		return "", err
	}
	fmt.Println(v, FormatUint(v), "++++")
	if toBase == 10 {
		return FormatUint(v), nil
	}
	return strconv.FormatUint(v, toBase), nil
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
	return strconv.FormatInt(v, 10)
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
func ToBool(s string) bool {
	v, _ := strconv.ParseBool(s)
	return v
}
func ToInt64(s string) int64 {
	v, _ := ParseInt64(s)
	return v
}
func ToInt(s string) int {
	v, _ := ParseInt(s)
	return v
}
func ToInt32(s string) int32 {
	v, _ := ParseInt32(s)
	return v
}
func ToInt24(s string) Int24 {
	v, _ := ParseInt24(s)
	return v
}
func ToInt16(s string) int16 {
	v, _ := ParseInt16(s)
	return v
}
func ToInt8(s string) int8 {
	v, _ := ParseInt8(s)
	return v
}
func ToUint64(s string) uint64 {
	v, _ := ParseUint64(s)
	return v
}
func ToUint(s string) uint {
	v, _ := ParseUint(s)
	return v
}
func ToUint32(s string) uint32 {
	v, _ := ParseUint32(s)
	return v
}
func ToUint24(s string) Uint24 {
	v, _ := ParseUint24(s)
	return v
}
func ToUint16(s string) uint16 {
	v, _ := ParseUint16(s)
	return v
}
func ToUint8(s string) uint8 {
	v, _ := ParseUint8(s)
	return v
}

func ParseInt(s string) (int, error) {
	return parseInt[int](s, 10, 32)
}
func ParseInt8(s string) (int8, error) {
	return parseInt[int8](s, 10, 8)
}
func ParseInt16(s string) (int16, error) {
	return parseInt[int16](s, 10, 16)
}
func ParseInt24(s string) (Int24, error) {
	v, err := parseInt[int32](s, 10, 32)
	return Int24(v), err
}
func ParseInt32(s string) (int32, error) {
	return parseInt[int32](s, 10, 32)
}
func ParseInt64(s string) (int64, error) {
	return parseInt[int64](s, 10, 64)
}

func ParseUint(s string) (uint, error) {
	return parseUint[uint](s, 10, 32)
}
func ParseUint8(s string) (uint8, error) {
	return parseUint[uint8](s, 10, 8)
}
func ParseUint16(s string) (uint16, error) {
	return parseUint[uint16](s, 10, 16)
}
func ParseUint24(s string) (Uint24, error) {
	v, err := parseUint[uint32](s, 10, 16)
	return Uint24(v), err
}
func ParseUint32(s string) (uint32, error) {
	return parseUint[uint32](s, 10, 32)
}
func ParseUint64(s string) (uint64, error) {
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
func parseInt[T ~int8 | ~int16 | ~int32 | ~int | ~int64](s string, base int, bitSize int) (T, error) {
	if num, ok := smallNumbersMap[s]; ok {
		return T(num), nil
	}
	v, err := strconv.ParseInt(s, base, bitSize)
	return T(v), err
}

// parseUint 通用无符号数字解析函数
// @example parseUint[uint8]("123", 10, 8)    parseUint[uint16]("123", 10, 8)
func parseUint[T ~uint8 | ~uint16 | ~uint32 | ~uint | ~uint64](s string, base int, bitSize int) (T, error) {
	if num, ok := smallNumbersMap[s]; ok {
		return T(num), nil
	}
	v, err := strconv.ParseUint(s, base, bitSize)
	return T(v), err
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
