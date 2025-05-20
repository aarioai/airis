package atype

import (
	"bytes"
	"encoding/binary"
	"github.com/aarioai/airis/pkg/types"
	"io"
	"strings"
)

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

/*
一般point 需要建 spatial 索引，就需要单独到一个表里，不应该放在一起
*/
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
