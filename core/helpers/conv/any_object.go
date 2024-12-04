package conv

import "github.com/aarioai/airis/core/atype"

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
