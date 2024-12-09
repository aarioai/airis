package arrmap

import (
	"github.com/aarioai/airis/pkg/types"
)

func MergeSlices[T types.BasicType](sources ...map[string][]T) map[string][]T {
	if len(sources) == 0 {
		return nil
	}
	if len(sources) == 1 {
		return sources[0]
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
	result := make(map[string][]T, totalSize)

	// 按顺序合并所有 map
	for _, m := range sources {
		if m == nil {
			continue
		}
		for k, v := range m {
			if len(v) == 0 {
				result[k] = nil
				continue
			}
			// 深拷贝切片
			newSlice := make([]T, len(v))
			copy(newSlice, v)
			result[k] = newSlice
		}
	}

	return result
}

// Merge 合并多个相同类型的 map，后面的 map 会覆盖前面的值
// 支持的类型包括基本类型，[]byte 应该用 MergeSlices
// 新开辟内存空间，可以避免副作用
func Merge[T types.BasicType](sources ...map[string]T) map[string]T {
	if len(sources) == 0 {
		return nil
	}
	if len(sources) == 1 {
		return sources[0]
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
	result := make(map[string]T, totalSize)
	for _, m := range sources {
		if m == nil {
			continue
		}
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
