package types

import (
	"cmp"
	"golang.org/x/exp/constraints"
	"reflect"
)

// Number is every numeric type.
type Number interface {
	constraints.Integer | constraints.Float
}

// BasicType is every basic type.
type BasicType interface {
	~bool | ~byte | ~rune | ~string | Number
}

// MapKeyType is an allowed map key type.
// Go map keys can also be structs, but that is not recommended.
type MapKeyType interface {
	cmp.Ordered // rune = int32, byte = uint8
}

// Stringable is a type that can be converted with string(T). The word stringable is used by Laravel, Apache, and others.
// ~ denotes the underlying type. For example, the underlying type of atype.Uint24 is uint32.
type Stringable interface {
	~[]byte | ~[]rune | ~byte | ~rune | ~string // rune = int32, byte = uint8
}

// For the other integer limits, use math.MaxInt64.
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

// Ref return the reference to the variable
func Ref[T any](n T) *T {
	return &n
}

// RefNz return the reference to the variable, return nil when the variable is zero
func RefNz[T ~byte | ~rune | ~string | Number](n T) *T {
	var zero T
	if n == zero {
		return nil
	}
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

// IsNil reports whether x is nil. In Go, an any value holds two pointers: one to the compile-time type, and one to the actual value.
func IsNil(x any) bool {
	// @warn Type assertions and reflection are not particularly fast. Use them only when necessary; limiting their use helps performance.
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

// IsUnsigned reports whether the type is an unsigned integer.
func IsUnsigned[T constraints.Integer](_ T) bool {
	var zero T
	switch any(zero).(type) {
	case uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}
