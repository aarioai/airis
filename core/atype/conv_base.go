package atype

import (
	"strconv"
	"sync/atomic"
)

// 为了提高性能，可以添加常用数字的字符串缓存。小数字在query string 或body传参中会很常见
var (
	smallNumbers [100]string
	initialized  uint32
)

func init() {
	if atomic.CompareAndSwapUint32(&initialized, 0, 1) {
		for i := 0; i < 100; i++ {
			smallNumbers[i] = strconv.Itoa(i)
		}
	}
}

// getSmallNumberString 获取小数字的字符串表示
func getSmallNumberString(n int) string {
	if n >= 0 && n < 100 {
		return smallNumbers[n]
	}
	return strconv.Itoa(n)
}
