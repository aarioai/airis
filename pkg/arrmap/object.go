package arrmap

import (
	"sort"
	"strings"
)

// JoinKeys 将map的key用sep连接起来
func JoinKeys[T any](m map[string]T, sep string, sorts ...bool) string {
	if len(m) == 0 {
		return ""
	}
	keys := Keys(m)
	if len(sorts) > 0 && sorts[0] {
		sort.Strings(keys)
	}

	return strings.Join(keys, sep)
}

// SortedKeys 获取map的key列表，并排序
func SortedKeys[T any](m map[string]T) []string {
	if len(m) == 0 {
		return nil
	}
	keys := Keys(m)
	sort.Strings(keys)
	return keys
}

// Keys 获取map的key列表
func Keys[T any](m map[string]T) []string {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
