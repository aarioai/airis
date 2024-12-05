package conv

import "github.com/aarioai/airis/core/atype"

// Deref 是一个通用的指针解引用函数，如果指针为 nil 则返回类型的零值
// T 必须是基本数值类型: uint64, uint32, uint, uint16, uint8, int64, int32, int, int16, int8, float64, float32, string
// @example a:=100; var b *int; Deref(&a)  Deref(b)
func Deref[T uint64 | uint32 | uint | uint16 | uint8 | int64 | int32 | int | int16 | int8 | float64 | float32 | string](n *T) T {
	if n == nil {
		var zero T
		return zero
	}
	return *n
}

// DerefInt24 专门处理 Int24 类型的指针解引用
func DerefInt24(n *atype.Uint24) atype.Uint24 {
	if n == nil {
		return 0
	}
	return *n
}

// DerefUint24 专门处理 Uint24 类型的指针解引用
func DerefUint24(n *atype.Uint24) atype.Uint24 {
	if n == nil {
		return 0
	}
	return *n
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

func AnyFloat64Map(mi map[string]any) (map[string]float64, error) {
	if len(mi) == 0 {
		return nil, nil
	}
	result := make(map[string]float64, len(mi))
	for k, v := range mi {
		val, err := atype.Float64(v, 64)
		if err != nil {
			return nil, err
		}
		result[k] = val
	}
	return result, nil
}

func AnyStrings(ai []any) []string {
	if len(ai) == 0 {
		return nil
	}
	result := make([]string, len(ai))
	for i, v := range ai {
		result[i] = atype.String(v)
	}
	return result
}

func AnyStringMap(mi map[string]any) map[string]string {
	if len(mi) == 0 {
		return nil
	}
	result := make(map[string]string, len(mi))
	for k, v := range mi {
		result[k] = atype.String(v)
	}
	return result
}

func AnyStringsMap(mi map[string]any) map[string][]string {
	if len(mi) == 0 {
		return nil
	}
	result := make(map[string][]string, len(mi))
	for k, v := range mi {
		if slice, ok := v.([]any); ok {
			result[k] = AnyStrings(slice)
		}
	}
	return result
}
func AnyComplexStringMap(mi map[string]any) map[string]map[string]string {
	if len(mi) == 0 {
		return nil
	}
	result := make(map[string]map[string]string, len(mi))
	for k, v := range mi {
		if m, ok := v.(map[string]any); ok {
			result[k] = AnyStringMap(m)
		}
	}
	return result
}
func AnyComplexStringsMap(mi map[string]any) map[string][][]string {
	if len(mi) == 0 {
		return nil
	}
	result := make(map[string][][]string, len(mi))
	for k, v := range mi {
		if slices, ok := v.([]any); ok {
			innerSlices := make([][]string, len(slices))
			for i, slice := range slices {
				if innerSlice, ok := slice.([]any); ok {
					innerSlices[i] = AnyStrings(innerSlice)
				}
			}
			result[k] = innerSlices
		}
	}
	return result
}
func AnyStringMaps(ai []any) []map[string]string {
	if len(ai) == 0 {
		return nil
	}
	var result []map[string]string
	for _, v := range ai {
		if m, ok := v.(map[string]any); ok {
			result = append(result, AnyStringMap(m))
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func AnyComplexMaps(ai []any) []map[string]any {
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
