package atype

import (
	"reflect"
)

// ToSlice 将任意值转换为 []any
func ToSlice(v any) []any {
	if v == nil {
		return nil
	}

	if slice, ok := v.([]any); ok {
		return slice
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil
	}

	result := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return result
}

// ToAnySlice 将任意类型切片转换为 []any
func ToAnySlice[T any](v []T) []any {
	if len(v) == 0 {
		return nil
	}
	args := make([]any, len(v))
	for i, x := range v {
		args[i] = x
	}
	return args
}

func ToStrings(ai []any) []string {
	if len(ai) == 0 {
		return nil
	}
	result := make([]string, len(ai))
	for i, v := range ai {
		result[i] = String(v)
	}
	return result
}
