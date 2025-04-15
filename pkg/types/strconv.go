package types

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"strconv"

	"golang.org/x/exp/constraints"
)

// 标准以：+/  ； URL 模式以 -_   ；  这里要贴合数字，尽量不要用数学符号，因此选用 _ ~
const Base64Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_~"

const NSmalls = math.MaxUint8 + 1 // small number less than this

// 为了提高性能，可以添加常用数字的字符串缓存。小数字在query string 或body传参中会很常见
// 虽然 strconv.ParseInt 里面是用数字字符串段方式逐步截取的。不过没有对常用小数缓存优化。这里对常用小数字用查表法
var (
	smallNumbers    [NSmalls]string // 暂时保持 uint8 最高范围
	smallNumberKeys = make(map[string]int, len(smallNumbers))

	smallBase64Numbers    [NSmalls]string
	smallBase64NumberKeys = make(map[string]int, len(smallBase64Numbers))
)

func init() {
	for i := 0; i < NSmalls; i++ {
		s := strconv.Itoa(i)
		smallNumbers[i] = s
		smallNumberKeys[s] = i

		b := formatBase64Uint(uint64(i))
		smallBase64Numbers[i] = b
		smallBase64NumberKeys[b] = i
	}
}
func ParseBase64Int(s string) (int64, error) {
	if s == "" {
		return 0, errors.New("empty string")
	}
	neg := s[0] == '-'
	if neg {
		s = s[1:]
	}
	n, err := ParseBase64Uint(s)
	if err != nil {
		return 0, err
	}
	if (!neg && n > math.MaxInt64) || n > math.MaxInt64+1 {
		return 0, fmt.Errorf("base64 int(%s) out of range: (-)%d", s, n)
	}
	if neg {
		n = -n
	}
	return int64(n), nil
}
func ParseBase64Uint(s string) (uint64, error) {
	return parseBase64Uint(s)
}
func parseBase64Uint(s string) (uint64, error) {
	if u, ok := smallBase64NumberKeys[s]; ok {
		return uint64(u), nil
	}
	base := 64
	c62 := Base64Digits[len(Base64Digits)-2]
	c63 := Base64Digits[len(Base64Digits)-1]
	cutoff := uint64(math.MaxUint64)/uint64(base) + 1
	var n uint64
	for _, c := range []byte(s) {
		var d byte
		switch {
		case c == c62:
			d = 62
		case c == c63:
			d = 63
		case '0' <= c && c <= '9':
			d = c - '0'
		case 'a' <= c && c <= 'z':
			d = c - 'a' + 10
		case 'A' <= c && c <= 'Z':
			d = c - 'A' + 36
		default:
			return 0, fmt.Errorf("invalid base64 digit: %s", string(c))
		}
		if n >= cutoff {
			// n*base overflows
			return math.MaxUint64, fmt.Errorf("base64 number %s is too large", s)
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > math.MaxUint64 {
			// n+d overflows
			return math.MaxUint64, fmt.Errorf("base64 number %s is too large", s)
		}
		n = n1
	}
	return n, nil
}
func FormatBase64Int[T constraints.Signed](u T) string {
	neg := u < 0
	if neg {
		u = -u
	}
	s := FormatBase64Uint(uint64(u))
	if neg {
		return "-" + s
	}
	return s
}
func FormatBase64Uint[T constraints.Unsigned](v T) string {
	u := uint64(v) // 必须要转换，否则无法比较
	if u < NSmalls {
		return smallBase64Numbers[int(u)]
	}
	return formatBase64Uint(u)
}

// 64 = 2^6 可以使用位移取代 %
func formatBase64Uint[T constraints.Unsigned](v T) string {
	base := uint(64)
	u := uint64(v)
	var a [64 + 1]byte // +1 for sign of 64bit value in base 2
	i := len(a)
	shift := uint(bits.TrailingZeros(base)) & 7
	m := uint64(base - 1) // == 1<<shift - 1
	for u > m {
		i--
		a[i] = Base64Digits[uint(u&m)]
		u >>= shift
	}
	// u < base
	i--
	a[i] = Base64Digits[uint(u)]
	return string(a[i:])
}

// 人为崩溃原则。strconv 也是这样做的
func checkBase(base int) {
	if base < 2 || (base > 32 && base != 64) {
		panic("atype: illegal number base")
	}
}

// ConvertBase 数字型字符串进制转换
// base 支持2~32，64
func ConvertBase(s string, fromBase, toBase int) (string, error) {
	if s == "" {
		return "", errors.New("string is empty")
	}
	if s == "0" {
		return "0", nil
	}
	checkBase(fromBase)
	checkBase(toBase)

	var err error
	if s[0] == '-' {
		var v int64
		switch fromBase {
		case 10:
			v, err = ParseInt64(s)
		case 64:
			v, err = ParseBase64Int(s)
		default:
			v, err = strconv.ParseInt(s, fromBase, 64)
		}
		if err != nil {
			return "", err
		}
		switch toBase {
		case 10:
			return FormatInt(v), nil
		case 64:
			return FormatBase64Int(v), nil
		default:
			return strconv.FormatInt(v, toBase), nil
		}
	}
	// unsigned
	var u uint64
	switch fromBase {
	case 10:
		u, err = ParseUint64(s)
	case 64:
		u, err = ParseBase64Uint(s)
	default:
		u, err = strconv.ParseUint(s, fromBase, 64)
	}
	if err != nil {
		return "", err
	}
	switch toBase {
	case 10:
		return FormatUint(u), nil
	case 64:
		return FormatBase64Uint(u), nil
	default:
		return strconv.FormatUint(u, toBase), nil
	}
}

// FormatUint8 直接查表
// 小数字在query string 或body传参中会很常见
func FormatUint8(v uint8) string {
	return smallNumbers[v]
}

// FormatInt 将int64转换为字符串，小数字直接使用查表法
func FormatInt[T constraints.Signed](value T) string {
	if value == 0 {
		return "0"
	}
	neg := value < 0
	if neg {
		value = -value
	}
	s := FormatUint(uint64(value))
	if neg {
		return "-" + s
	}
	return s
}
func Itoa[T constraints.Signed](value T) string {
	return FormatInt(value)
}

// FormatUint 将uint64转换为字符串，小数字直接使用查表法
func FormatUint[T constraints.Unsigned](value T) string {
	v := uint64(value)
	// 小数字使用查表法
	if v < NSmalls {
		return smallNumbers[v]
	}
	return strconv.FormatUint(v, 10)
}

// FormatFloat 将float转换为字符串，整数部分使用FormatInt，小数部分使用strconv.FormatFloat
func FormatFloat[T constraints.Float](value T, bitSize ...int) string {
	v := float64(value)
	intPart := int64(v)
	// 处理int性能更好
	if v == float64(intPart) {
		return FormatInt(intPart)
	}
	bs := 64
	if len(bitSize) > 0 {
		bs = bitSize[0]
	}
	return strconv.FormatFloat(v, 'f', -1, bs)
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

func ToUint16(s string) uint16 {
	v, _ := ParseUint16(s)
	return v
}
func ToUint8(s string) uint8 {
	v, _ := ParseUint8(s)
	return v
}
func ToFloat32(s string) float32 {
	result, _ := strconv.ParseFloat(s, 32)
	return float32(result)
}
func ToFloat64(s string) float64 {
	result, _ := strconv.ParseFloat(s, 64)
	return result
}
func Atoi(s string) (int, error) {
	return ParseInt(s)
}
func ParseInt(s string) (int, error) {
	return ParseSigned[int](s, 10, 32)
}
func ParseInt8(s string) (int8, error) {
	return ParseSigned[int8](s, 10, 8)
}
func ParseInt16(s string) (int16, error) {
	return ParseSigned[int16](s, 10, 16)
}

func ParseInt32(s string) (int32, error) {
	return ParseSigned[int32](s, 10, 32)
}
func ParseInt64(s string, bitSize ...int) (int64, error) {
	if len(bitSize) > 1 {
		panic("ParseInt64 invalid bit size")
	}
	bs := 64
	if len(bitSize) > 0 {
		bs = bitSize[0]
	}
	return ParseSigned[int64](s, 10, bs)
}

func ParseUint(s string) (uint, error) {
	return ParseUnsigned[uint](s, 10, 32)
}
func ParseUint8(s string) (uint8, error) {
	return ParseUnsigned[uint8](s, 10, 8)
}
func ParseUint16(s string) (uint16, error) {
	return ParseUnsigned[uint16](s, 10, 16)
}

func ParseUint32(s string) (uint32, error) {
	return ParseUnsigned[uint32](s, 10, 32)
}
func ParseUint64(s string, bitSize ...int) (uint64, error) {
	if len(bitSize) > 1 {
		panic("ParseUint64 invalid bit size")
	}
	bs := 64
	if len(bitSize) > 0 {
		bs = bitSize[0]
	}
	return ParseUnsigned[uint64](s, 10, bs)
}

// ParseSigned 通用数字解析函数
// @example ParseSigned[int8]("123", 10, 8)    ParseSigned[uint8]("123", 10, 8)
func ParseSigned[T ~int8 | ~int16 | ~int32 | ~int | ~int64](s string, base int, bitSize int) (T, error) {
	if num, ok := smallNumberKeys[s]; ok {
		return T(num), nil
	}
	v, err := strconv.ParseInt(s, base, bitSize)
	return T(v), err
}

// ParseUnsigned 通用无符号数字解析函数
// @example ParseUnsigned[uint8]("123", 10, 8)    ParseUnsigned[uint16]("123", 10, 8)
func ParseUnsigned[T ~uint8 | ~uint16 | ~uint32 | ~uint | ~uint64](s string, base int, bitSize int) (T, error) {
	if num, ok := smallNumberKeys[s]; ok {
		return T(num), nil
	}
	v, err := strconv.ParseUint(s, base, bitSize)
	return T(v), err
}
