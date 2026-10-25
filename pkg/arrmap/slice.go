package arrmap

import "golang.org/x/exp/slices"

// slices.Compact([]) []  removes consecutive duplicates --> [5,1,2,2,3,3,4,5,5] ==> [5, 1,2,3,4,5] len=6(), cap()=9

// slices.Contains([], value:Any)  bool  whether an element is in the slice
// slices.ContainsFunc([], func(i int)bool{})
// slices.BinarySearch([], value:Any) index, bool   reports whether an element is in the slice and returns its index
// slices.Index([], value:Any) int
// slices.Equal([], [])  bool  corresponding elements of the two slices must match exactly
// slices.Compare([], [])  0 when corresponding elements are identical; -1 when the element at the same position in the first slice is lexicographically greater than that in the second, otherwise 1
// slices.Delete([], start, end) []  ===>  append(a[0:start], a[end:]...)  exactly equivalent
// slices.Insert([], start, values...) []   inserts elements starting at a position
// slices.ReplaceAll([], start, end, values...)  replaces elements within a range
// slices.IsSorted([]) bool
// slices.Sort([]) []
// slices.Revers([]) []  reverses
// slices.Max([]) / slices.Min([])

// Performance
// slices.Clip([])  drops unused capacity from the slice
// slices.Clone([]) []  shallow-copies the slice
// slices.Grow([], n)

func clearSlice[S ~[]E, E any](s S) {
	var zero E
	for i := range s {
		s[i] = zero
	}
}

// Compact remove duplicated slice items
// slices.Compact() only removes consecutive duplicates, so sort first
// slices.Compact([]) []  removes consecutive duplicates --> [5,1,2,2,3,3,4,5,5] ==> [5, 1,2,3,4,5] len=6(), cap()=9
// This deduplicates the whole slice.
func Compact[S ~[]E, E comparable](s S, sorted bool) S {
	if sorted {
		return slices.Compact(s)
	}
	if len(s) < 2 {
		return s
	}

	i := 1
	for k := 1; k < len(s); k++ {
		var exists bool
		for j := 0; j < i; j++ {
			if s[j] == s[k] {
				exists = true
				break
			}
		}
		if !exists {
			if i != k {
				s[i] = s[k]
			}
			i++
		}
	}
	clearSlice(s[i:]) // zero/nil out the obsolete elements, for GC
	return s[:i]
}

func HasDuplicates[T comparable](arr []T) bool {
	seen := make(map[T]bool)
	for _, item := range arr {
		if seen[item] {
			return true
		}
		seen[item] = true
	}
	return false
}
