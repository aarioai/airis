package atype

import (
	"errors"
	"fmt"
	"strconv"
)

// String convert into string
// @warn byte is a built-in alias of uint8, String('A') returns "97"; String(AByte('A')) returns "A"
func String(d any) string {
	if d == nil {
		return ""
	}
	switch v := d.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	// Name(Abyte('A')) returns "A"
	case AByte:
		return string([]byte{byte(v)})
	case Date:
		return string(v)
	case Datetime:
		return string(v)
	case bool:
		return strconv.FormatBool(v)

		// 整数
	case int64:
		return FormatInt(v)
	case int:
		return FormatInt(int64(v))
	case int32:
		return FormatInt(int64(v))
	//  is a built-in alias of int32, @notice 'A' is a rune(65), is different with byte('A') (alias of uint8(65))
	//case rune: return FormatInt(int64(v))
	case Int24:
		return FormatInt(int64(v))
	case int16:
		return FormatInt(int64(v))
	case int8:
		return FormatInt(int64(v))

		// 无符号整数
	case uint64:
		return FormatUint(v)
	case uint:
		return FormatUint(uint64(v))
	case uint32:
		return FormatUint(uint64(v))
	case Uint24:
		return FormatUint(uint64(v))
	case uint16:
		return FormatUint(uint64(v))
	case uint8: // byte
		return FormatUint(uint64(v))
	// is a built-in alias of uint8, Name('A') returns "97"
	//case byte: return strconv.FormatUint(uint64(v), 10)
	case Booln:
		return FormatUint(uint64(v))

		// 浮点数
	case float64:
		return FormatFloat(v, 64)
	case float32:
		return FormatFloat(float64(v), 32)
	}
	// 有些类型type vt uint  var a, b vt 这样就无法识别为 uint；所以尝试通过字符串方式转一下
	return fmt.Sprint(d)
}

func Bytes(d any) []byte {
	if d == nil {
		return nil
	}
	switch v := d.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	case AByte:
		return []byte{byte(v)}
	}
	// 很少会有number/bool转bytes的情况，因此，不需要过度优化
	// 其他类型通过 String 转换
	return []byte(String(d))
}

func Bool(d any) (bool, error) {
	if d == nil {
		return false, nil
	}

	switch v := d.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case Booln:
		return v.Bool(), nil
	case int8:
		return v > 0, nil
	case uint8:
		return v > 0, nil
	case int:
		return v > 0, nil
	}
	// 有些类型type vt uint  var a, b vt 这样就无法识别为 uint；所以尝试通过字符串方式转一下
	return strconv.ParseBool(String(d))
}

func Int8(d any) (int8, error) {
	v, err := Int64Base(d, 8)
	return int8(v), err
}

func Int16(d any) (int16, error) {
	v, err := Int64Base(d, 16)
	return int16(v), err
}
func Int32(d any) (int32, error) {
	v, err := Int64Base(d, 32)
	return int32(v), err
}

func Int(d any) (int, error) {
	v, err := Int64Base(d, 32)
	return int(v), err
}
func Int64(d any) (int64, error) {
	v, err := Int64Base(d, 64)
	return v, err
}
func Int64Base(d any, bitSize int) (int64, error) {
	if d == nil {
		return 0, errors.New("nil to integer")
	}

	switch v := d.(type) {
	case AByte:
		return int64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseInt(v, 10, bitSize)
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case rune:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case byte:
		return int64(v), nil
	case Booln:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case Uint24:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	}
	// 有些类型type vt uint  var a, b vt 这样就无法识别为 uint；所以尝试通过字符串方式转一下
	return strconv.ParseInt(String(d), 10, bitSize)
}

func Uint8(d any) (uint8, error) {
	v, err := Uint64Base(d, 8)
	return uint8(v), err
}

func Uint16(d any) (uint16, error) {
	v, err := Uint64Base(d, 16)
	return uint16(v), err
}
func Uint24b(d any) (Uint24, error) {
	v, err := Uint64Base(d, 24)
	return Uint24(v), err
}
func Uint32(d any) (uint32, error) {
	r, err := Uint64Base(d, 32)
	return uint32(r), err
}

func Uint(d any) (uint, error) {
	r, err := Uint64Base(d, 32)
	return uint(r), err
}
func Uint64(d any) (uint64, error) {
	return Uint64Base(d, 64)
}
func Uint64Base(d any, bitSize int) (uint64, error) {
	if d == nil {
		return 0, errors.New("nil to uint64")
	}
	switch v := d.(type) {
	case AByte:
		return uint64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseUint(v, 10, bitSize)
	case int8:
		return uint64(v), nil
	case int16:
		return uint64(v), nil
	case rune:
		return uint64(v), nil
	case int:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case byte: // 等同uint8
		return uint64(v), nil
	case Booln:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case Uint24:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint64:
		return v, nil
	case float32:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	}
	// 有些类型type vt uint  var a, b vt 这样就无法识别为 uint；所以尝试通过字符串方式转一下
	return strconv.ParseUint(String(d), 10, bitSize)
}
func Float32(d any) (float32, error) {
	f, err := Float64(d, 32)
	if err != nil {
		return 0.0, err
	}
	return float32(f), nil
}
func Float64(d any, bitSize int) (float64, error) {
	if d == nil {
		return 0.0, errors.New("nil to float64")
	}
	switch v := d.(type) {
	case AByte:
		return float64(v), nil
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil

	case string:
		return strconv.ParseFloat(v, bitSize)
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case rune:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case byte:
		return float64(v), nil
	case Booln:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case Uint24:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	}
	// 有些类型type vt uint  var a, b vt 这样就无法识别为 uint；所以尝试通过字符串方式转一下
	return strconv.ParseFloat(String(d), bitSize)
}
