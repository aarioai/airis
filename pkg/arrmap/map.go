package arrmap

import (
	"github.com/aarioai/airis/pkg/types"
)

// maps.Copy(dst, src)
// maps.Clone(m)
// maps.Equal(a, b)

// Merge 将多个map合并为一个新的map
// @warn 这里是浅拷贝
func Merge[A types.MapKeyType, T any](sources ...map[A]T) map[A]T {
	if len(sources) == 0 {
		return nil
	}

	// 计算总容量以预分配空间
	totalSize := 0
	for _, m := range sources {
		totalSize += len(m)
	}
	if totalSize == 0 {
		return nil
	}
	// 创建结果 map，预分配合适的容量
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
