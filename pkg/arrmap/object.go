package arrmap

import (
	"github.com/aarioai/airis/pkg/types"
	"slices"
	"strings"
)

// JoinKeys joins the map keys with sep.
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

// SortedKeys returns the map keys, sorted.
func SortedKeys[A types.MapKeyType, T any](m map[A]T) []A {
	keys := Keys(m)
	if len(keys) == 0 {
		return nil
	}
	slices.Sort(keys)
	return keys
}

// SortedKeysFunc returns the map keys, sorted. The filter rule is similar to bytes.TrimFunc().
func SortedKeysFunc[A types.MapKeyType, T any](m map[A]T, filter func(key A, value T) (A, bool)) []A {
	keys := KeysFunc(m, filter)
	if len(keys) == 0 {
		return nil
	}
	slices.Sort(keys)
	return keys
}

// Keys returns the map keys.
// maps.Keys() is an iter.Seq[K], used as: for k := range maps.Keys(m)
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

// KeysFunc returns the map keys. The filter rule is similar to bytes.TrimFunc().
// maps.Keys() is an iter.Seq[K], used as: for k := range maps.Keys(m)
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
