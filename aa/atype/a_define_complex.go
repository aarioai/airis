package atype

type Bitwise struct {
	BitName  string // 该位名称
	BitPos   BitPos // big endian 下，位所在位置
	BitValue bool   // 该位的值
	MaxBits  uint8
}
type Bitwiser struct {
	BitName  string      // 该位名称
	BitPos   BitPosition // big endian 下，位所在位置
	BitValue bool        // 该位的值
	MaxBits  uint8
}
