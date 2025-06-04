package arrmap

import "golang.org/x/exp/slices"

// slices.Compact([]) []  移除连续重复项 --> [5,1,2,2,3,3,4,5,5] ==> [5, 1,2,3,4,5] len=6(), cap()=9

// slices.Contains([], value:Any)  bool 某个元素是否在slice里
// slices.ContainsFunc([], func(i int)bool{})
// slices.BinarySearch([], value:Any) index, bool   检测某个元素是否在slice里面，同时返回索引
// slices.Index([], value:Any) int
// slices.Equal([], [])  bool 两个切片位置对应元素要完全一致
// slices.Compare([], [])  若位置对应元素完全相同，0； -1 第一slice相同位置元素字母大于第二个slice相同位置元素字母，反之 1
// slices.Delete([], start, end) []  ===>  append(a[0:start], a[end:]...)  完全等同
// slices.Insert([], start, values...) []   在某个位置开始，插入一些元素
// slices.ReplaceAll([], start, end, values...)  在区间内替换
// slices.IsSorted([]) bool
// slices.Sort([]) []
// slices.Revers([]) []  反转
// slices.Max([]) / slices.Min([])

// 性能
// slices.Clip([])  删除切片中未使用的容量
// slices.Clone([]) [] 浅拷贝slice
// slices.Grow([], n)

func clearSlice[S ~[]E, E any](s S) {
	var zero E
	for i := range s {
		s[i] = zero
	}
}

// Compact remove duplicated slice items
// slices.Compact() 仅对连续的重复去重，因此最好先进行 sort
// slices.Compact([]) []  移除连续重复项 --> [5,1,2,2,3,3,4,5,5] ==> [5, 1,2,3,4,5] len=6(), cap()=9
// 这里对全部去重
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
