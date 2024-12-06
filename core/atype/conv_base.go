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
func fastIntToString(v int64) string {
	// 小数字使用查表法；小数字在query string 或body传参中会很常见
	if v >= 0 && v < 100 {
		return getSmallNumberString(int(v))
	}
	return strconv.FormatInt(v, 10)
}

func fastUintToString(v uint64) string {
	// 小数字使用查表法
	if v < 100 {
		return getSmallNumberString(int(v))
	}

	return strconv.FormatUint(v, 10)
}

func fastFloatToString(v float64, bitSize int) string {
	// 整数部分处理
	if v == float64(int64(v)) {
		return fastIntToString(int64(v))
	}

	// 使用 strconv.FormatFloat
	return strconv.FormatFloat(v, 'f', -1, bitSize)
}
