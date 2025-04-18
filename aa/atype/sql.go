package atype

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"errors"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
	"io"
	"net"
	"strings"
	"time"
)

type NullUint64 struct{ sql.NullInt64 }
type NullString struct{ sql.NullString }

type Location struct {
	Valid     bool    `json:"-"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Height    float64 `json:"height"` // 保留
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}
type Coordinate struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
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
type Ip struct{ sql.NullString } //  VARBINARY(16) | BINARY(16) 固定16位长度 net.IP               // IP Address

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
type Bin string // binary string
type Booln uint8
type Int24 int32
type Uint24 uint32
type Year uint16      // uint16 date: yyyy
type YearMonth Uint24 // uint24 date: yyyymm  不要用 Date，主要是不需要显示dd。
type YMD uint         // YYYYMMDD
type Date string      // yyyy-mm-dd
type Datetime string  // yyyy-mm-dd hh:ii:ss
type UnixTime int64   // int 形式 datetime，可与 datetime, date 互转

type SepStrings string // a,b,c,d,e
type SepUint8s string  // 1,2,3,4
type SepUint16s string // 1,2,3,4
type SepUint24s string // 1,2,3,4
type SepUint32s string // 1,2,3,4
type SepInts string    // 1,2,3,4
type SepUints string   // 1,2,3,4
type SepUint64s string // 1,2,3,4

const (
	False       Booln    = 0
	True        Booln    = 1
	MinDate     Date     = "0000-00-00"
	MaxDate     Date     = "9999-12-31"
	MinDatetime Datetime = "0000-00-00 00:00:00"
	MaxDatetime Datetime = "9999-12-31 23:59:59"
)

func NewNullUint64(value uint64) NullUint64 {
	var v NullUint64
	if value > 0 {
		v.Scan(value)
	}
	return v
}
func (t NullUint64) Uint64() uint64 {
	if !t.Valid {
		return 0
	}
	return uint64(t.Int64)
}
func (t NullUint64) Equal(b uint64) bool {
	if !t.Valid {
		return false
	}
	return t.Uint64() == b
}
func NewNullString(value string) NullString {
	var v NullString
	if value != "" {
		v.Scan(value)
	}
	return v
}

func (t NullString) Equal(b string) bool {
	if !t.Valid {
		return false
	}
	return t.String == b
}

func delimiter(delimiters ...string) string {
	if len(delimiters) == 0 || delimiters[0] == "" {
		return ","
	}
	return delimiters[0]
}
func (l Location) Coordinate() *Coordinate {
	if !l.Valid {
		return nil
	}
	return &Coordinate{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Height:    l.Height,
	}
}
func (l Location) Position() Position {
	return ToPosition(l.Coordinate())
}

// https://dev.mysql.com/doc/refman/8.0/en/gis-data-formats.html
//	The value length is 25 bytes, made up of these components (as can be seen from the hexadecimal value):
//	4 bytes for integer SRID (0)       4326 是GPS   WGS84，表示按 lat-lng 保存
//	1 byte for integer byte order (1 = little-endian)
//	4 bytes for integer type information (1 = Point)
//	8 bytes for double-precision X coordinate (1)
//	8 bytes for double-precision Y coordinate (−1)

func ToPositionBase(srid uint32, order byte, typ uint32, x, y float64) Position {
	var pos Position
	buf := new(bytes.Buffer)
	buf.Grow(25)
	// uint32就是4个字节
	binary.Write(buf, binary.LittleEndian, srid)
	binary.Write(buf, binary.LittleEndian, order)
	binary.Write(buf, binary.LittleEndian, typ)
	binary.Write(buf, binary.LittleEndian, x)
	binary.Write(buf, binary.LittleEndian, y)
	pos.Scan(buf.Bytes())
	return pos
}
func ToPosition(coord *Coordinate) Position {
	if coord == nil {
		return Position{}
	}
	//  4326 是GPS   WGS84，表示按 lat-lng 保存
	return ToPositionBase(4326, 1, 1, coord.Latitude, coord.Longitude)
}
func (p Position) Bytes() []byte {
	if !p.Ok() {
		return nil
	}
	return []byte(p.String)
}
func (p Position) Ok() bool {
	return p.Valid && len(p.String) == 25
}
func binaryRead(r io.Reader, littleEndian bool, data any) error {
	if littleEndian {
		return binary.Read(r, binary.LittleEndian, data)
	}
	return binary.Read(r, binary.BigEndian, data)
}

func (p Position) Parse() (srid uint32, order byte, typ uint32, x float64, y float64, ok bool) {
	b := p.Bytes()
	if b == nil {
		return
	}
	buf := bytes.NewReader(b[4:5])
	binary.Read(buf, binary.LittleEndian, &order) // 只有1字节，无论bigEndian，还是littleEndian，结果都一样
	littleEndian := order == 1
	buf = bytes.NewReader(b[0:4])
	binaryRead(buf, littleEndian, &srid)
	buf = bytes.NewReader(b[5:9])
	binaryRead(buf, littleEndian, &typ)
	buf = bytes.NewReader(b[9:17])
	binaryRead(buf, littleEndian, &x)
	buf = bytes.NewReader(b[17:25])
	binaryRead(buf, littleEndian, &y)
	return
}
func (p Position) Coordinate() *Coordinate {
	_, _, _, x, y, ok := p.Parse()
	if !ok {
		return nil
	}
	return &Coordinate{
		Latitude:  x,
		Longitude: y,
	}
}
func (p Position) Point() *Point {
	_, _, _, x, y, ok := p.Parse()
	if !ok {
		return nil
	}
	return &Point{
		X: x,
		Y: y,
	}
}

// sql 可以使用 select *， cast(ip as CHAR) from table_name   进行显示
func ToIp(addr string) Ip {
	var ip Ip
	if addr == "" {
		return ip
	}
	nip := net.ParseIP(addr) // 无论是IPv4还是IPv6都是16字节
	if nip == nil {
		return ip
	}
	ip2 := nip.To4() // 将IPv4的转为4字节
	if ip2 != nil {
		nip = ip2
	}

	ip.Scan(nip.String())
	return ip
}
func (ip Ip) Bytes() []byte {
	if !ip.Ok() {
		return nil
	}
	return []byte(ip.String)
}

func (ip Ip) Ok() bool {
	return ip.Valid && len(ip.String) == net.IPv4len || len(ip.String) == net.IPv6len
}
func (ip Ip) Net() net.IP {
	b := ip.Bytes()
	if b == nil {
		return nil
	}
	return b
}

// 是不是 IPv4
func (ip Ip) Is4() bool {
	return len(ip.String) == net.IPv4len
}

// 无论是4字节的IPv4，还是16字节的IPv4或IPv6，都能输出可阅读的IP地址
func (ip Ip) To16() string {
	ip2 := ip.Net()
	// 包括ipv4 / ipv16
	nip := ip2.To16() // 此时IP长度为16
	if nip == nil {
		return ""
	}
	return nip.String() // 返回IPv4样式IP地址
}

func (n Uint24) Uint32() uint32 { return uint32(n) }

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
func (b Bin) StmtValue() string {
	return "b'" + string(b) + "'"
}

func ToBooln(b bool) Booln {
	if b {
		return True
	}
	return False
}
func NewBooln(b uint8) Booln {
	if b > 0 {
		return True
	}
	return False
}
func (b Booln) Uint8() uint8 {
	if b.Bool() {
		return 1
	}
	return 0
}
func (b Booln) Bool() bool    { return b > 0 }
func (b Booln) IsFalse() bool { return b == 0 }
func (b Booln) IsTrue() bool  { return b > 0 }
func ToYearMonth(year int, month time.Month) YearMonth {
	if year < 0 {
		return 0
	}
	ym := year*100 + int(month)
	return YearMonth(ym)
}

func ParseYear(year string) (Year, error) {
	y, err := types.ParseUint16(year)
	if err != nil {
		return 0, err
	}
	return Year(y), nil
}
func ParseYearMonth(year, month string) (YearMonth, error) {
	y, err1 := types.ParseUint16(year)
	m, err2 := types.ParseUint8(month)
	if err := ae.FirstError(err1, err2); err != nil {
		return 0, err
	}
	return YearMonth(uint32(y)*100 + uint32(m)), nil
}
func (n Int24) String() string     { return types.FormatInt(n) }
func (n Uint24) String() string    { return types.FormatUint(n) }
func (y Year) String() string      { return types.FormatUint(y) }
func (y YearMonth) String() string { return types.FormatUint(y) }
func (y YMD) String() string       { return types.FormatUint(y) }

// years/months 可为负数
func (y YearMonth) Add(years int, months int, loc *time.Location) YearMonth {
	tm := y.Time(loc).AddDate(years, months, 0)
	return ToYearMonth(tm.Year(), tm.Month())
}
func (y YearMonth) Uin32() uint32 { return uint32(y) }
func (y YearMonth) Date() (int, time.Month) {
	year := int(y) / 100
	month := time.Month(y % 100)
	return year, month
}
func (y YearMonth) Time(loc *time.Location) time.Time {
	year, month := y.Date()
	return time.Date(year, month, 0, 00, 00, 00, 0, loc)
}

func ParseYMD(year, month, day string) (YMD, error) {
	y, err1 := types.ParseUint16(year)
	m, err2 := types.ParseUint8(month)
	d, err3 := types.ParseUint8(day)
	if err := ae.FirstError(err1, err2, err3); err != nil {
		return 0, err
	}
	return YMD(uint(y)*10000 + uint(m)*100 + uint(d)), nil
}
func (y YMD) Uint() uint {
	return uint(y)
}
func (y YMD) Date() Date {
	s := types.FormatUint(y)
	if len(s) != 8 {
		return MinDate
	}
	x := s[0:4] + "-" + s[4:6] + "-" + s[6:]
	return Date(x)
}

// time.Now().In()  loc 直接通过 in 传递
func NewDate(d string, loc *time.Location) Date {
	if d == "" || d == MinDate.String() {
		return MinDate
	}
	_, err := time.ParseInLocation("2006-01-02", d, loc)
	if err != nil {
		return MinDate
	}
	return Date(d)
}
func ToDate(t *time.Time) Date {
	if t == nil || t.IsZero() {
		return MinDate
	}
	return Date(t.Format("2006-01-02"))
}
func (d Date) Valid() bool {
	return len(d) == 10 && d != MinDate && d != MaxDate && d != "1970-01-01"
}
func (d Date) String() string { return string(d) }
func (d Date) Time(loc *time.Location) (time.Time, error) {
	if !d.Valid() {
		return time.Time{}, errors.New("min date")
	}
	return time.ParseInLocation("2006-01-02", string(d), loc)
}
func (d Date) OrMin() Date {
	if !d.Valid() {
		return MinDate
	}
	return d
}
func (d Date) OrMax() Date {
	if !d.Valid() {
		return MaxDate
	}
	return d
}
func (d Date) OrNow(loc *time.Location) Date {
	if !d.Valid() {
		return Date(time.Now().In(loc).Format("2006-01-02"))
	}
	return d
}

func (d Date) Ymd() YMD {
	s := strings.ReplaceAll(d.String(), "-", "")
	n, _ := types.ParseUint(s)
	return YMD(n)
}
func (d Date) Int64(loc *time.Location) int64 {
	if !d.Valid() {
		return 0
	}
	tm, err := time.ParseInLocation("2006-01-02", string(d), loc)
	if err != nil {
		return 0
	}
	return tm.Unix()
}
func (d Date) Unix(loc *time.Location) UnixTime {
	if !d.Valid() {
		return 0
	}
	t, _ := d.Time(loc)
	return UnixTime(t.Unix())
}

func NewDatetime(d string, loc *time.Location) Datetime {
	if d == "" || d == MinDatetime.String() {
		return MinDatetime
	}
	_, err := time.ParseInLocation("2006-01-02 15:04:05", d, loc)
	if err != nil {
		return MinDatetime
	}
	return Datetime(d)
}
func Now(loc *time.Location) Datetime {
	return ToDatetime(time.Now(), loc)
}

func ToDatetime(t time.Time, loc *time.Location) Datetime {
	if t.IsZero() {
		return MinDatetime
	}
	return Datetime(t.In(loc).Format("2006-01-02 15:04:05"))
}
func ToDatetime2(t *time.Time, loc *time.Location) Datetime {
	if t == nil || t.IsZero() {
		return MinDatetime
	}
	return ToDatetime(*t, loc)
}
func UnixToDatetime(u int64, loc *time.Location) Datetime { return NewUnixTime(u).Datetime(loc) }

func (d Datetime) Valid() bool {
	return len(d) == 19 && d != MinDatetime && d != MaxDatetime && d.String() != "1970-01-01 00:00:00"
}
func (d Datetime) String() string { return string(d) }
func (d Datetime) Time(loc *time.Location) (time.Time, error) {
	if !d.Valid() {
		return time.Time{}, errors.New("min datetime")
	}
	return time.ParseInLocation("2006-01-02 15:04:05", string(d), loc)
}
func (d Datetime) Date() Date { return Date(d[0:len(MinDate)]) }
func (d Datetime) OrMin() Datetime {
	if !d.Valid() {
		return MinDatetime
	}
	return d
}
func (d Datetime) OrMax() Datetime {
	if !d.Valid() {
		return MaxDatetime
	}
	return d
}
func (d Datetime) OrNow(loc *time.Location) Datetime {
	if !d.Valid() {
		return Now(loc)
	}
	return d
}
func (d Datetime) Int64(loc *time.Location) int64 {
	if !d.Valid() {
		return 0
	}
	tm, err := time.ParseInLocation("2006-01-02 15:04:05", string(d), loc)
	if err != nil {
		return 0
	}
	return tm.Unix()
}
func (d Datetime) Unix(loc *time.Location) UnixTime {
	if !d.Valid() {
		return 0
	}
	t, _ := d.Time(loc)
	return UnixTime(t.Unix())
}

func NewUnixTime(u int64) UnixTime {
	return UnixTime(u)
}
func (u UnixTime) Int64() int64 { return int64(u) }
func (u UnixTime) Date(loc *time.Location) Date {
	if u == 0 {
		return MinDate
	}
	t := time.Unix(u.Int64(), 0).In(loc)
	return ToDate(&t)
}
func (u UnixTime) Datetime(loc *time.Location) Datetime {
	if u == 0 {
		return MinDatetime
	}
	t := time.Unix(u.Int64(), 0)
	return ToDatetime(t, loc)
}

func ToSepStrings(elems []string, delimiters ...string) SepStrings {
	return SepStrings(strings.Join(elems, delimiter(delimiters...)))
}
func (t SepStrings) Strings(delimiters ...string) []string {
	if t == "" {
		return nil
	}
	return strings.Split(string(t), delimiter(delimiters...))
}

func ToSepUint8s(elems []uint8, delimiters ...string) SepUint8s {
	// strings.Concat 类同
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return SepUint8s(types.FormatUint(elems[0]))
	}
	deli := delimiter(delimiters...)
	n := len(deli)*(len(elems)-1) + (len(elems) * types.MaxUint8Len)

	var b strings.Builder
	b.Grow(n)
	b.WriteString(types.FormatUint(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(deli)
		b.WriteString(types.FormatUint(s))
	}

	return SepUint8s(b.String())
}
func (t SepUint8s) Uint8s(delimiters ...string) []uint8 {
	if t == "" {
		return nil
	}
	arr := strings.Split(string(t), delimiter(delimiters...))
	v := make([]uint8, len(arr))
	for i, a := range arr {
		v[i] = types.ToUint8(a)
	}
	return v
}

func ToSepUint16s(elems []uint16, delimiters ...string) SepUint16s {
	// strings.Concat 类同
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return SepUint16s(types.FormatUint(elems[0]))
	}
	deli := delimiter(delimiters...)
	n := len(deli)*(len(elems)-1) + (len(elems) * types.MaxUint16Len)
	var b strings.Builder
	b.Grow(n)
	b.WriteString(types.FormatUint(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(deli)
		b.WriteString(types.FormatUint(s))
	}

	return SepUint16s(b.String())
}
func (t SepUint16s) Uint16s(delimiters ...string) []uint16 {
	if t == "" {
		return nil
	}
	arr := strings.Split(string(t), delimiter(delimiters...))
	v := make([]uint16, len(arr))
	for i, a := range arr {
		v[i] = types.ToUint16(a)
	}
	return v
}

func ToSepUint24s(elems []Uint24, delimiters ...string) SepUint24s {
	// strings.Concat 类同
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return SepUint24s(types.FormatUint(elems[0]))
	}
	deli := delimiter(delimiters...)
	n := len(deli)*(len(elems)-1) + (len(elems) * types.MaxUint24Len)
	var b strings.Builder
	b.Grow(n)
	b.WriteString(types.FormatUint(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(deli)
		b.WriteString(types.FormatUint(s))
	}

	return SepUint24s(b.String())
}

func (t SepUint24s) Uint32s(delimiters ...string) []Uint24 {
	if t == "" {
		return nil
	}
	arr := strings.Split(string(t), delimiter(delimiters...))
	v := make([]Uint24, len(arr))
	for i, a := range arr {
		v[i] = ToUint24(a)
	}
	return v
}

func ToSepInts(elems []int, delimiters ...string) SepInts {
	// strings.Concat 类同
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return SepInts(types.FormatInt(elems[0]))
	}
	deli := delimiter(delimiters...)
	n := len(deli)*(len(elems)-1) + (len(elems) * types.MaxIntLen)
	var b strings.Builder
	b.Grow(n)
	b.WriteString(types.FormatInt(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(deli)
		b.WriteString(types.FormatInt(s))
	}

	return SepInts(b.String())
}
func (t SepInts) Ints(delimiters ...string) []int {
	if t == "" {
		return nil
	}
	arr := strings.Split(string(t), delimiter(delimiters...))
	v := make([]int, len(arr))
	for i, a := range arr {
		b, _ := types.ParseInt(a)
		v[i] = int(b)
	}
	return v
}

func ToSepUints(elems []uint, delimiters ...string) SepUints {
	// strings.Concat 类同
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return SepUints(types.FormatUint(elems[0]))
	}
	deli := delimiter(delimiters...)
	n := len(deli)*(len(elems)-1) + (len(elems) * types.MaxUintLen)
	var b strings.Builder
	b.Grow(n)
	b.WriteString(types.FormatUint(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(deli)
		b.WriteString(types.FormatUint(s))
	}

	return SepUints(b.String())
}
func (t SepUints) Uints(delimiters ...string) []uint {
	if t == "" {
		return nil
	}
	arr := strings.Split(string(t), delimiter(delimiters...))
	v := make([]uint, len(arr))
	for i, a := range arr {
		x, _ := types.ParseUint(a)
		v[i] = x
	}
	return v
}

func ToSepUint64s(elems []uint64, delimiters ...string) SepUint64s {
	// strings.Concat 类同
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return SepUint64s(types.FormatUint(elems[0]))
	}
	deli := delimiter(delimiters...)
	n := len(deli)*(len(elems)-1) + (len(elems) * types.MaxUint64Len)
	var b strings.Builder
	b.Grow(n)
	b.WriteString(types.FormatUint(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(deli)
		b.WriteString(types.FormatUint(s))
	}

	return SepUint64s(b.String())
}
func (t SepUint64s) Uint64s(delimiters ...string) []uint64 {
	if t == "" {
		return nil
	}
	arr := strings.Split(string(t), delimiter(delimiters...))
	v := make([]uint64, len(arr))
	for i, a := range arr {
		v[i], _ = types.ParseUint64(a)
	}
	return v
}
