package atype

import (
	"github.com/aarioai/airis/pkg/types"
	"reflect"
)

// ToMap 将任意值转换为 map[any]any
func ToMap(value any) (map[any]any, bool) {
	switch v := value.(type) {
	case map[any]any:
		return v, true
	case map[string]any:
		m := make(map[any]any, len(v))
		for k, val := range v {
			m[k] = val
		}
		return m, true
	case map[string]string:
		m := make(map[any]any, len(v))
		for k, val := range v {
			m[k] = val
		}
		return m, true
	}

	// 使用反射处理其他 map 类型
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Map {
		return nil, false
	}

	m := make(map[any]any, rv.Len())
	for _, key := range rv.MapKeys() {
		m[key.Interface()] = rv.MapIndex(key).Interface()
	}
	return m, true
}
func ToStringMaps(ai []any) []map[string]string {
	if len(ai) == 0 {
		return nil
	}
	var result []map[string]string
	for _, v := range ai {
		if m, ok := v.(map[string]any); ok {
			result = append(result, ToStringMap(m))
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func ToComplexMaps(ai []any) []map[string]any {
	if len(ai) == 0 {
		return nil
	}
	var result []map[string]any
	for _, v := range ai {
		if m, ok := v.(map[string]any); ok {
			result = append(result, m)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func ToFloat64Map[T types.MapKeyType](mi map[T]any) (map[T]float64, error) {
	if len(mi) == 0 {
		return nil, nil
	}
	result := make(map[T]float64, len(mi))
	for k, v := range mi {
		val, err := Float64(v, 64)
		if err != nil {
			return nil, err
		}
		result[k] = val
	}
	return result, nil
}

func ToStringMap[T types.MapKeyType](mi map[T]any) map[T]string {
	if len(mi) == 0 {
		return nil
	}
	result := make(map[T]string, len(mi))
	for k, v := range mi {
		result[k] = String(v)
	}
	return result
}

func ToStringsMap[T types.MapKeyType](mi map[T]any) map[T][]string {
	if len(mi) == 0 {
		return nil
	}
	result := make(map[T][]string, len(mi))
	for k, v := range mi {
		if slice, ok := v.([]any); ok {
			result[k] = ToStrings(slice)
		}
	}
	return result
}
func ToComplexStringMap[T types.MapKeyType](mi map[T]any) map[T]map[string]string {
	if len(mi) == 0 {
		return nil
	}
	result := make(map[T]map[string]string, len(mi))
	for k, v := range mi {
		if m, ok := v.(map[string]any); ok {
			result[k] = ToStringMap(m)
		}
	}
	return result
}
func ToComplexStringsMap[T types.MapKeyType](mi map[T]any) map[T][][]string {
	if len(mi) == 0 {
		return nil
	}
	result := make(map[T][][]string, len(mi))
	for k, v := range mi {
		if slices, ok := v.([]any); ok {
			innerSlices := make([][]string, len(slices))
			for i, slice := range slices {
				if innerSlice, ok := slice.([]any); ok {
					innerSlices[i] = ToStrings(innerSlice)
				}
			}
			result[k] = innerSlices
		}
	}
	return result
}
