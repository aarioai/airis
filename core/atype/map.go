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
	if m.Value == nil {
		return nil, fmt.Errorf("map is nil")
	}
	current, ok := ToMap(m.Value)
	if !ok {
		return nil, fmt.Errorf("invalid map: %v", m.Value)
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

	// 遍历所有键
	for i := 0; i < len(allKeys)-1; i++ {
		nextValue, exists := current[allKeys[i]]
		if !exists {
			return nil, fmt.Errorf("key not found: %s", strings.Join(ToStrings(allKeys[:i+1]), "."))
		}

		nextMap, ok := ToMap(nextValue)
		if !ok {
			return nil, fmt.Errorf("invalid intermediate value at key %s: %v", strings.Join(ToStrings(allKeys[:i+1]), "."), nextValue)
		}
		current = nextMap
	}

	// 获取最终值
	finalKey := allKeys[len(allKeys)-1]
	finalValue, exists := current[finalKey]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", strings.Join(ToStrings(allKeys), "."))
	}

	return finalValue, nil
}
