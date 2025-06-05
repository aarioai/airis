package atype

import (
	"github.com/aarioai/airis/pkg/basic"
	"github.com/aarioai/airis/pkg/types"
	"time"
)

func (s Sid) Nullable() NullableSid {
	result := NullableSid{}
	if s != "" {
		result.Scan(string(s))
	}
	return result
}

func NewBooln(n uint8) (Booln, bool) { return basic.Ter(n == 1, True, False), n == 0 || n == 1 }
func ToBooln(b bool) Booln           { return basic.Ter(b, True, False) }
func (b Booln) Uint8() uint8         { return basic.Ter(b == True, uint8(1), uint8(0)) }
func (b Booln) Valid() bool          { return b == True || b == False }
func (b Booln) Bool() bool           { return b > 0 }
func (b Booln) IsFalse() bool        { return b == False }
func (b Booln) IsTrue() bool         { return b == True }

func NewChar(b byte) (Char, bool) { return Char(b), b > 31 && b < 127 }
func (c Char) String() string     { return string(c) }

func NewBin(s string) (Bin, bool) { return Bin(s), IsBin(s) }
func (b Bin) Normalize() string   { return "b'" + string(b) + "'" }

// Uint8 BitPos https://en.wikipedia.org/wiki/Bit_numbering
func (b BitPos) Uint8() uint8        { return uint8(b) }
func (b BitPosition) Uint16() uint16 { return uint16(b) }

// SET x=x|v
func (b Bitwise) SetStmt(fieldName string) string {
	if b.BitValue {
		bv := 1 << b.BitPos
		bs := types.FormatInt(bv)
		return fieldName + "=" + fieldName + "|" + bs
	}
	return b.unsetStmt(fieldName)
}

func (b Bitwise) unsetStmt(fieldName string) string {
	maxBits := (1 << b.MaxBits) - 1
	bv := maxBits - (1 << b.BitPos)
	bs := types.FormatInt(bv)
	return fieldName + "=" + fieldName + "&" + bs
}

func ToMillisecond(t time.Duration) Millisecond { return Millisecond(t / time.Millisecond) }
func (s Millisecond) Duration() time.Duration   { return time.Duration(s) * time.Millisecond }
func ToSecond(t time.Duration) Second           { return Second(t / time.Second) }
func (s Second) Duration() time.Duration        { return time.Duration(s) * time.Second }
func (s Second) Add(t time.Duration) Second     { return ToSecond(s.Duration() + t) }

func (n Uint24) Uint32() uint32 { return uint32(n) }

func (n Int24) String() string  { return types.FormatInt(n) }
func (n Uint24) String() string { return types.FormatUint(n) }

func ToInt24(s string) Int24 {
	v, _ := ParseInt24(s)
	return v
}
func ToUint24(s string) Uint24 {
	v, _ := ParseUint24(s)
	return v
}
func ParseInt24(s string) (Int24, error) {
	v, err := types.ParseSigned[int32](s, 10, 24)
	return Int24(v), err
}
func ParseUint24(s string) (Uint24, error) {
	v, err := types.ParseUnsigned[uint32](s, 10, 24)
	return Uint24(v), err
}

// DerefInt24 专门处理 Int24 类型的指针解引用
func DerefInt24(n *Uint24) Uint24 {
	if n == nil {
		return 0
	}
	return *n
}

// DerefUint24 专门处理 Uint24 类型的指针解引用
func DerefUint24(n *Uint24) Uint24 {
	if n == nil {
		return 0
	}
	return *n
}
