package types

import "golang.org/x/exp/constraints"

// Number 表示所有数字类型
type Number interface {
	constraints.Integer | constraints.Float
}

// BasicType 表示所有基本类型
type BasicType interface {
	bool | byte | string | Number
}

// 其他的使用 math.MaxInt64
const MaxInt24 = 1<<23 - 1
const MinInt24 = -1 << 23
const MaxUint24 = 1<<24 - 1

const MaxInt8Len = 4    // -128 ~ 127
const MaxUint8Len = 3   // 0 ~ 256
const MaxInt16Len = 6   // -32768 ~ 32767
const MaxUint16Len = 5  // 0 ~ 65535
const MaxInt24Len = 8   // -8388608 ~ 8388607
const MaxUint24Len = 8  // 0 ~ 16777215
const MaxIntLen = 11    // -2147483648 ~ 2147483647
const MaxUintLen = 10   // 0 ~ 4294967295
const MaxInt64Len = 20  // -9223372036854775808 ~ 9223372036854775807
const MaxUint64Len = 20 // 0 ~ 18446744073709551615
