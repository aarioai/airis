package arrmap

import (
	"github.com/aarioai/airis/pkg/types"
)

// maps.Copy(dst, src)
// maps.Clone(m)
// maps.Equal(a, b)

// Merge combines multiple maps into a new map.
// @warn This is a shallow copy.
func Merge[A types.MapKeyType, T any](sources ...map[A]T) map[A]T {
	if len(sources) == 0 {
		return nil
	}

	// Sum the sizes so the result can be preallocated.
	totalSize := 0
	for _, m := range sources {
		totalSize += len(m)
	}
	if totalSize == 0 {
		return nil
	}
	// Create the result map with a suitable preallocated capacity.
	result := make(map[A]T, totalSize)
	for _, m := range sources {
		for k, v := range m {
			result[k] = v
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
