package atype

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

// 其他的使用 math.MaxInt64
const MaxInt24 = 1<<23 - 1
const MinInt24 = -1 << 23
const MaxUint24 = 1<<24 - 1

const MaxInt8Len = 4    // -128 ~ 127
const MaxUint8Len = 3   // 0 ~ 256
const MaxInt16Len = 6   // -32768 ~ 32767
const MaxUint16Len = 5  // 0 ~ 65535
const MaxInt24Len = 8   // -8388608 ~ 8388607
const MaxUint24Len = 8  // 0 ~ 16777215
const MaxIntLen = 11    // -2147483648 ~ 2147483647
const MaxUintLen = 10   // 0 ~ 4294967295
const MaxInt64Len = 20  // -9223372036854775808 ~ 9223372036854775807
const MaxUint64Len = 20 // 0 ~ 18446744073709551615

const MaxDistriLen = MaxUint24Len

// Invalid Kind = iota
// Bool
// Int
// Int8
// Int16
// Int32
// Int64
// Uint
// Uint8
// Uint16
// Uint32
// Uint64
// Uintptr
// Float32
// Float64
// Complex64
// Complex128
// Array
// Chan
// Func
// Interface
// Map
// Ptr
// Slice
// String
// Struct
// UnsafePointer
// 获取原始类型  i 用指针
// @param i 必须为指针
// @return 除了 reflect.Ptr 外其他类型；包括 interface
func PrimitiveType(i any) reflect.Kind {
	if i == nil {
		return reflect.Invalid // nil
	}
	k := reflect.TypeOf(i).Elem().Kind()
	if k == reflect.Invalid {
		return reflect.Invalid // nil
	}
	if k == reflect.Ptr {
		v := reflect.ValueOf(i).Elem()
		if !v.CanInterface() {
			return reflect.Invalid
		}
		return PrimitiveType(v.Interface())
	}
	if k == reflect.Interface {
		k = reflect.ValueOf(i).Kind()
		if k == reflect.Ptr {
			v := reflect.ValueOf(i).Elem()
			if !v.CanInterface() {
				return reflect.Invalid
			}
			return PrimitiveType(v.Interface())
		}
		return k
	}
	return k
}

// 可能为指针，或者其他
func PType(i any) reflect.Kind {
	if i == nil {
		return reflect.UnsafePointer // nil
	}
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	k := t.Kind()
	if k == reflect.Invalid {
		return reflect.Invalid // nil
	}
	// 指针
	if k == reflect.Ptr {
		return PrimitiveType(i)
	}
	if k == reflect.Interface {
		k = reflect.ValueOf(i).Kind()
		if k == reflect.Ptr {
			v = reflect.ValueOf(i).Elem()
			if !v.CanInterface() {
				return reflect.Invalid
			}
			return PrimitiveType(v.Interface())
		}
		return k
	}
	return k
}
func IsEmpty(d any) bool {
	return !NotEmpty(d)
}

/*
在Go语言中，一个any类型的变量包含了2个指针，一个指针指向值的在编译时确定的类型，另外一个指针指向实际的值。
*/
func IsNil(x any) bool {
	// @warn 断言和反射性能不是特别好，如果不得已再使用，控制使用有助于提升程序性能。
	return x == nil || (reflect.ValueOf(x).Kind() == reflect.Ptr && reflect.ValueOf(x).IsNil())
}

// NotEmpty check if a value is empty
// @warn NotEmpty(byte(0)) == false,  NotEmpty(byte('0')) == true
//
//	NotEmpty(0) == false, NotEmpty(-1) == true, NotEmpty(1) == true
func NotEmpty(d any) bool {
	if d == nil {
		return false
	}
	switch v := d.(type) {
	case Abyte:
		return v != 0
	case bool:
		return v
	case string:
		return v != ""
	case int8:
		return v != 0 // 复数，不算 empty，只有0 才算empty
	case int16:
		return v != 0
	case rune:
		return v != 0
	case int:
		return v != 0
	case int64:
		return v != 0
	case byte:
		return v > 0
	case Booln:
		return v > 0
	case uint16:
		return v > 0
	case Uint24:
		return v > 0
	case uint32:
		return v > 0
	case uint:
		return v > 0
	case uint64:
		return v > 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	}

	return String(d) != ""
}

/*
"Array and slice values encode as JSON arrays, except that []byte encodes as a base64-encoded string, and a nil slice encodes as the null JSON object.
json.Marshal() 不能正常转换 []byte 及 []uint8
*/
func MarshalUint8s(x []uint8) ([]byte, error) {
	if x == nil {
		return nil, nil
	}
	var b bytes.Buffer
	b.WriteByte('[')
	for i, c := range x {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatUint(uint64(c), 10))
	}
	b.WriteByte(']')
	return b.Bytes(), nil
}
func UnmarshalUint8s(x []byte) ([]uint8, error) {
	if x == nil || len(x) < 2 {
		return nil, nil
	}

	n := len(x) - 1 // remove last ']'
	v := make([]uint8, 0)
	for i := 1; i < n; {
		for x[i] == ' ' || x[i] == ',' {
			i++
		}
		if x[i] < '0' || x[i] > '9' {
			return nil, errors.New("invalid uint8 json: " + string(x))
		}
		var c string
		for x[i] >= '0' && x[i] <= '9' {
			c += string(x[i])
			i++
		}
		u, _ := strconv.ParseUint(c, 10, 8)
		v = append(v, uint8(u))
	}
	return v, nil
}
