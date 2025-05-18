package atype

import "database/sql"

type Bin string  // binary string
type Booln uint8 // 0 | 1

type Int24 int32
type Uint24 uint32

type Decimal int64 // [ -922337203685477.5808,  -922337203685477.5807]
type Money Decimal // 有效范围：正负100亿元；  ±100 0000亿
type VMoney Money  // 1 coin = 1 money    如 chatgpt 等消耗，单次消耗低于0.1分，因此需要更大的 coin比例

type RoundType uint8

type Year uint16      // 4 digits number, format YYYY
type YearMonth Uint24 // 6 digits number, format YYYYMM
type YMD uint         // 8 digits number, format YYYYMMDD
type Date string      // format YYYY-MM-DD
type Datetime string  // format YYYY-MM-DD HH:II:SS
type Timestamp int64  // unix timestamp

// type Html template.HTML   HTML 直接使用 template.HTML
type Province uint8  //2 digits province district code
type Dist uint16     //  4 digits district code
type Distri Uint24   // 6 digits district code
type District uint64 // 12 digits district code

// Version 版本，nginx 等比较版本，就是转为每节3位数字的整数进行比较
// Semantic Versioning https://semver.org/lang/zh-CN/
// [tag 0 release, 1 alpha, 2 beta,3 RC, 4 Revision][major 000-999][minor 000-999][build/patch 000-999]
type Version uint
type VersionTag uint8
type VersionStruct struct {
	Main  uint // Major*1000000 + Minor*1000 + Patch
	Major uint
	Minor uint
	Patch uint
	Tag   VersionTag
}

type NullUint64 struct{ sql.NullInt64 }
type NullString struct{ sql.NullString }

type Location struct {
	Valid     bool    `json:"-"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Height    float64 `json:"height"` // 保留
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Height    float64 `json:"height"` // 保留
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

/*
一般point 需要建 spatial 索引，就需要单独到一个表里，不应该放在一起
*/
type Position struct{ sql.NullString } // []byte // postion, coordinate or point

// sql 可以使用 select *， cast(ip as CHAR) from table_name   进行显示
type IP struct{ sql.NullString } //  VARBINARY(16) | BINARY(16) 固定16位长度 net.IP               // IP Address

// https://en.wikipedia.org/wiki/Bit_numbering
type BitPos uint8       // bit-position (in big endian)
type BitPosition uint16 // bit-position (in big endian)
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

type SepStrings string // a,b,c,d,e
type SepUint8s string  // 1,2,3,4
type SepUint16s string // 1,2,3,4
type SepUint24s string // 1,2,3,4
type SepUint32s string // 1,2,3,4
type SepInts string    // 1,2,3,4
type SepUints string   // 1,2,3,4
type SepUint64s string // 1,2,3,4

type Text string // Text 65535 bytes

type File string
type Document string
type Image string
type Video string
type Audio string
type Files struct{ NullStrings }
type Documents struct{ NullStrings }
type Images struct{ NullStrings }
type Videos struct{ NullStrings }
type Audios struct{ NullStrings }
