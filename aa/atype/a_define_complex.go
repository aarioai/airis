package atype

type Bitwise struct {
	BitName  string // name of this bit
	BitPos   BitPos // bit position under big-endian
	BitValue bool   // value of this bit
	MaxBits  uint8
}
type Bitwiser struct {
	BitName  string      // name of this bit
	BitPos   BitPosition // bit position under big-endian
	BitValue bool        // value of this bit
	MaxBits  uint8
}
