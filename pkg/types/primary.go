package types

import (
	"reflect"
)

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
// Get the primitive type. Pass i as a pointer.
// @param i must be a pointer
// @return any kind other than reflect.Ptr, including interface
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

// May be a pointer, or some other kind.
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
	// Pointer.
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
