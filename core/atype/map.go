package atype

import (
	"fmt"
	"strings"
)

// structs.Map(rsp) 可以将struct 转为 map[string]any
type Map struct {
	Value any
}

func NewMap(v any) Map {
	return Map{
		Value: v,
	}
}

// Get key from a map[string]any
// p.Get("users.1.name") 等同于 p.Get("user", "1", "name")
// @warn p.Get("user", "1", "name") 与 p.Get("user", 1, "name") 不一样
func (m Map) Get(key any, keys ...any) (any, error) {
	value := m.Value
	if value == nil {
		return nil, fmt.Errorf("map is nil")
	}
	val, ok := ToMap(value)
	if !ok {
		return nil, fmt.Errorf("invalid map: %v", value)
	}
	// 处理点号分隔的路径
	allKeys := make([]any, 0, len(keys)+1)
	if strKey, ok := key.(string); ok && strings.Contains(strKey, ".") {
		ks := ToAnySlice(strings.Split(strKey, "."))
		allKeys = append(allKeys, ks...)
	} else {
		allKeys = append(allKeys, key)
		allKeys = append(allKeys, keys...)
	}

	for i, k := range allKeys {
		value, ok = val[k]
		if !ok {
			return nil, fmt.Errorf("key not found: %s", strings.Join(ToStrings(allKeys), "."))
		}

		// 如果是最后一个键，直接返回值
		if i == len(allKeys)-1 {
			return value, nil
		}
	}

	return value, nil
}
