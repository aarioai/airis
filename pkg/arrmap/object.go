package arrmap

import (
	"github.com/aarioai/airis/pkg/types"
	"slices"
	"strings"
)

// JoinKeys 将map的key用sep连接起来
func JoinKeys[A byte | rune | string](m map[A]any, sep string, sort bool) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, string(k))
	}
	if sort {
		slices.Sort(keys)
	}
	return strings.Join(keys, sep)
}

// SortedKeys 获取map的key列表，并排序
func SortedKeys[A types.MapKeyType, T any](m map[A]T) []A {
	keys := Keys(m)
	if len(keys) == 0 {
		return nil
	}
	slices.Sort(keys)
	return keys
}

// SortedKeysFunc 获取map的key列表，并排序。命令规则类似 bytes.TrimFunc()
func SortedKeysFunc[A types.MapKeyType, T any](m map[A]T, filter func(key A, value T) (A, bool)) []A {
	keys := KeysFunc(m, filter)
	if len(keys) == 0 {
		return nil
	}
	slices.Sort(keys)
	return keys
}

// Keys 获取map的key列表
// maps.Keys() 是一个iter.Seq[K]，通过 for  k := range maps.Keys(m) 使用
func Keys[A types.MapKeyType, T any](m map[A]T) []A {
	if len(m) == 0 {
		return nil
	}
	keys := make([]A, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// KeysFunc 获取map的key列表。命令规则类似 bytes.TrimFunc()
// maps.Keys() 是一个iter.Seq[K]，通过 for  k := range maps.Keys(m) 使用
func KeysFunc[A types.MapKeyType, T any](m map[A]T, filter func(key A, value T) (A, bool)) []A {
	if len(m) == 0 {
		return nil
	}
	keys := make([]A, 0, len(m))
	for k := range m {
		if newK, ok := filter(k, m[k]); ok {
			keys = append(keys, newK)
		}
	}
	return keys
}
