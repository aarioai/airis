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

type VersionStruct struct {
	Main  uint // Major*1000000 + Minor*1000 + Patch
	Major uint
	Minor uint
	Patch uint
	Tag   VersionTag
}
