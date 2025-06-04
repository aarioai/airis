package arrmap

import "golang.org/x/exp/constraints"

func BuildSliceMap[K comparable, V constraints.Integer](s []K) map[K]V {
	m := make(map[K]V, len(s))
	for i, c := range s {
		m[c] = V(i)
	}
	return m
}
