package arrmap

// Extend 将后面的map合并到第一个map中，后面的map会覆盖前面的值
func Extend[T any](base map[string]T, sources ...map[string]T) {
	if base == nil {
		panic("base map is nil")
	}
	for _, m := range sources {
		for k, v := range m {
			base[k] = v
		}
	}
}

// Merge 将多个map合并为一个新的map
// @warn 这里是浅拷贝
func Merge[T any](sources ...map[string]T) map[string]T {
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
	result := make(map[string]T, totalSize)
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

// Copy 创建一个 map 的浅拷贝
// 如果源 map 为 nil，则返回 nil
func Copy[T any](source map[string]T) map[string]T {
	if source == nil {
		return nil
	}
	result := make(map[string]T, len(source))
	for k, v := range source {
		result[k] = v
	}
	return result
}
