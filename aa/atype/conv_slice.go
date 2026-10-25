package atype

import (
	"errors"
	"reflect"
)

// Slice converts any value to []any
func Slice(v any) ([]any, error) {
	if v == nil {
		return nil, nil
	}

	if slice, ok := v.([]any); ok {
		return slice, nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, errors.New("cast type error")
	}

	result := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return result, nil
}

// ToAnySlice converts a slice of any type to []any
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

func ToStrings[T any](ai []T) []string {
	if len(ai) == 0 {
		return nil
	}
	result := make([]string, len(ai))
	for i, v := range ai {
		result[i] = String(v)
	}
	return result
}
