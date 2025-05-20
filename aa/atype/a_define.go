package atype

import "database/sql"

// See ./t_basic.go

type Byte byte
type Int24 int32
type Uint24 uint32

type Bin string         // binary string
type BitPos uint8       // bit-position (in big endian)
type BitPosition uint16 // bit-position (in big endian)
type Booln uint8        // 0 | 1
type Millisecond int64
type Second int64

// See ./t_path_param.go

type PathParamType string // e.g. :string, :uuid, :uint8
type UUID string          // 32 or 36 bytes, 8-4-4-4-12
type Digits string        // \d+
type Lowers string        // [a-z]+
type Uppers string        // [A-Z]+
type Alphabetical string  // [a-zA-Z]+
type AlphaDigits string   // [a-zA-Z\d]+
type Word string          // \w+
type Email string
type Weekday uint8 // [0-6] from sunday to saturday

// See ./t_path.go
type FileType string        // e.g. .jpg
type StdFilename string     // [\w-.]+ // standard file name
type Filename string        // [\w-.!@#$%^&(){}~]+ , unicode filename
type StdPath string         // [\w-.\/]+
type Path string            // [\w-.!@#$%^&(){}~/]+ , unicode path
type URL string             // e.g. https://xxx/video.avi?quality=80
type FilenamePattern string //  e.g. {name}_{size:int}.avi
type PathPattern string     //  e.g. /a/b/{name}_{size:int}.avi
type UrlPattern string      // e.g. https://xxx/{user}/video.avi?quality={quality:int}

// See ./t_path_mime.go

type FilePath string
type DocumentPath FilePath
type ImagePath FilePath
type VideoPath FilePath
type AudioPath FilePath
type FilePaths struct{ NullStrings }
type DocumentPaths FilePaths
type ImagePaths FilePaths
type VideoPaths FilePaths
type AudioPaths FilePaths

// See ./t_decimal.go

type RoundType uint8
type Decimal int64 // [ -922337203685477.5808,  -922337203685477.5807]
type Money Decimal // 有效范围：正负100亿元；  ±100 0000亿
type VMoney Money  // 1 coin = 1 money    如 chatgpt 等消耗，单次消耗低于0.1分，因此需要更大的 coin比例

// See ./t_date.go

type Year uint16      // 4 digits number, format YYYY
type YearMonth Uint24 // 6 digits number, format YYYYMM
type YMD uint         // 8 digits number, format YYYYMMDD
type Date string      // format YYYY-MM-DD
type Datetime string  // format YYYY-MM-DD HH:mm:ss
type Timestamp int64  // unix timestamp

// See ./t_district.go

type Province uint8  //2 digits province district code
type Dist uint16     //  4 digits district code
type Distri Uint24   // 6 digits district code
type District uint64 // 12 digits district code

// See ./t_version.go

type Version uint
type VersionTag uint8

// See ./t_object.go

type NullUint64 struct{ sql.NullInt64 }
type NullString struct{ sql.NullString }
type NullJson struct{ sql.NullString }              // any
type NullUint8s struct{ sql.NullString }            // uint8 json array
type NullUint16s struct{ sql.NullString }           // uint16 json array
type NullUint24s struct{ sql.NullString }           // Uint24 json array
type NullUint32s struct{ sql.NullString }           // uint32 json array
type NullInts struct{ sql.NullString }              // int json array
type NullUints struct{ sql.NullString }             // uint json array
type NullUint64s struct{ sql.NullString }           // uint64 json array
type NullStrings struct{ sql.NullString }           // string json array
type NullStringMap struct{ sql.NullString }         // map[string]string   // JSON 规范，key 必须为字符串
type NullStringMapsMap struct{ sql.NullString }     // map[string][]map[string]string
type NullStringsMap struct{ sql.NullString }        // map[string][]string
type NullComplexStringMap struct{ sql.NullString }  // map[string]map[string]string
type NullComplexStringsMap struct{ sql.NullString } // map[string][][]string
type NullStringMaps struct{ sql.NullString }        // []map[string]string
type NullComplexStringMaps struct{ sql.NullString } // []map[string][]map[string]string
type ComplexMaps struct{ sql.NullString }           // []map[string]any

// See ./t_sep.go

type SepStrings string // a,b,c,d,e
type SepUint8s string  // 1,2,3,4
type SepUint16s string // 1,2,3,4
type SepUint24s string // 1,2,3,4
type SepUint32s string // 1,2,3,4
type SepInts string    // 1,2,3,4
type SepUints string   // 1,2,3,4
type SepUint64s string // 1,2,3,4

// See ./t_text.go

type Text string // Text 65535 bytes

// See ./t_complex_bytes.go

type Position struct{ sql.NullString } // []byte // postion, coordinate or point
type IP struct{ sql.NullString }       //  VARBINARY(16) | BINARY(16) 固定16位长度 net.IP               // IP Address
