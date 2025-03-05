package types

import (
	"cmp"
	"golang.org/x/exp/constraints"
	"reflect"
)

// Number 表示所有数字类型
type Number interface {
	constraints.Integer | constraints.Float
}

// BasicType 表示所有基本类型
type BasicType interface {
	~bool | ~byte | ~rune | ~string | Number
}

// MapKeyType 允许的map key类型
// 事实上Go map key 还支持struct，但是不建议使用
type MapKeyType interface {
	cmp.Ordered // rune = int32, byte = uint8
}

// Stringable 可以直接使用 string(T) 转换的类型。  单词 stringable 被Laravel、Apache等使用
// ~ 表示底层类型，如 atype.Uint24 底层类型是 uint32
type Stringable interface {
	~[]byte | ~[]rune | ~byte | ~rune | ~string // rune = int32, byte = uint8
}

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

// Ref returns the reference to the variable
func Ref[T any](n T) *T {
	return &n
}

// Deref dereference a pointer to its value. return zero on the pointer is nil
func Deref[T any](n *T) T {
	if n == nil {
		var zero T
		return zero
	}
	return *n
}

// IsNil 在Go语言中，一个any类型的变量包含了2个指针，一个指针指向值的在编译时确定的类型，另外一个指针指向实际的值。
func IsNil(x any) bool {
	// @warn 断言和反射性能不是特别好，如果不得已再使用，控制使用有助于提升程序性能。
	if x == nil {
		return true
	}
	switch reflect.ValueOf(x).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(x).IsNil()
	}
	return false
}

func IsEmpty(v any) bool {
	if IsNil(v) {
		return true
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Invalid:
		return true
	case reflect.Bool:
		return !v.(bool)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(v).Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(v).Uint() == 0
	case reflect.Uintptr:
		return reflect.ValueOf(v).Uint() == 0
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(v).Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return reflect.ValueOf(v).Complex() == 0
	case reflect.Array:
		return reflect.ValueOf(v).Len() == 0
	//case reflect.Chan:
	//case reflect.Func:
	//case reflect.Interface:
	case reflect.Map:
		return reflect.ValueOf(v).Len() == 0
	//case reflect.Pointer:
	case reflect.Slice:
		return reflect.ValueOf(v).Len() == 0
	case reflect.String:
		return reflect.ValueOf(v).Len() == 0
	//case reflect.Struct:
	//case reflect.UnsafePointer:
	//return reflect.ValueOf(v).IsNil()
	default:
		return false
	}
}
func NotEmpty(v any) bool {
	return !IsEmpty(v)
}

// IsUnsigned 判断类型是否为无符号整数
func IsUnsigned[T constraints.Integer](_ T) bool {
	var zero T
	switch any(zero).(type) {
	case uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}
