package atype

import (
	"encoding/json"
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

func NewNullJson(s []byte) NullJson {
	var x NullJson
	if len(s) > 0 {
		x.Scan(string(s))
	}
	return x
}

func ToNullJson(v any) NullJson {
	if v == nil {
		return NullJson{}
	}
	s, _ := json.Marshal(v)
	return NewNullJson(s)
}

func (t ComplexMaps) Interface() any {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v any
	json.Unmarshal([]byte(t.String), &v)
	return v
}

func NewNullUint8s(s string) NullUint8s {
	var x NullUint8s
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullUint8s(v []uint8) NullUint8s {
	if len(v) == 0 {
		return NullUint8s{}
	}
	s, _ := MarshalUint8s(v)
	if len(s) == 0 {
		return NullUint8s{}
	}

	return NewNullUint8s(string(s))
}

func (t NullUint8s) Uint8s() []uint8 {
	if !t.Valid || t.String == "" {
		return nil
	}
	w, _ := UnmarshalUint8s([]byte(t.String))
	return w
}
func NewNullUint16s(s string) NullUint16s {
	var x NullUint16s
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullUint16s(v []uint16) NullUint16s {
	if len(v) == 0 {
		return NullUint16s{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullUint16s{}
	}

	return NewNullUint16s(string(s))
}

func (t NullUint16s) Uint16s() []uint16 {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []any
	json.Unmarshal([]byte(t.String), &v)
	if len(v) == 0 {
		return nil
	}
	w := make([]uint16, len(v))
	newV := New()
	defer newV.Close()
	for i, x := range v {
		w[i] = newV.Reload(x).DefaultUint16(0)
	}
	return w
}
func NewNullUint24s(s string) NullUint24s {
	var x NullUint24s
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullUint24s(v []uint32) NullUint24s {
	if len(v) == 0 {
		return NullUint24s{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullUint24s{}
	}
	return NewNullUint24s(string(s))
}

func (t NullUint24s) Uint24s() []Uint24 {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []any
	json.Unmarshal([]byte(t.String), &v)
	if len(v) == 0 {
		return nil
	}
	w := make([]Uint24, len(v))
	newV := New()
	defer newV.Close()
	for i, x := range v {
		w[i] = newV.Reload(x).DefaultUint24(0)
	}
	return w
}
func NewNullUint32s(s string) NullUint32s {
	var x NullUint32s
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullUint32s(v []uint32) NullUint32s {
	if len(v) == 0 {
		return NullUint32s{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullUint32s{}
	}
	return NewNullUint32s(string(s))
}

func (t NullUint32s) Uint32s() []uint32 {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []any
	json.Unmarshal([]byte(t.String), &v)
	if len(v) == 0 {
		return nil
	}
	w := make([]uint32, len(v))
	newV := New()
	defer newV.Close()
	for i, x := range v {
		w[i] = newV.Reload(x).DefaultUint32(0)
	}
	return w
}
func NewNullInts(s string) NullInts {
	var x NullInts
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullInts(v []int) NullInts {
	if len(v) == 0 {
		return NullInts{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullInts{}
	}
	return NewNullInts(string(s))
}

func (t NullInts) Ints() []int {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []any
	json.Unmarshal([]byte(t.String), &v)
	if len(v) == 0 {
		return nil
	}
	w := make([]int, len(v))
	newV := New()
	defer newV.Close()
	for i, x := range v {
		w[i] = newV.Reload(x).DefaultInt(0)
	}
	return w
}
func NewNullUints(s string) NullUints {
	var x NullUints
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullUints(v []uint) NullUints {
	if len(v) == 0 {
		return NullUints{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullUints{}
	}

	return NewNullUints(string(s))
}
func (t NullUints) Uints() []uint {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []any
	json.Unmarshal([]byte(t.String), &v)
	if len(v) == 0 {
		return nil
	}
	w := make([]uint, len(v))
	newV := New()
	defer newV.Close()
	for i, x := range v {
		w[i] = newV.Reload(x).DefaultUint(0)
	}
	return w
}

func NewNullUint64s(s string) NullUint64s {
	var x NullUint64s
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullUint64s(v []uint64) NullUint64s {
	if len(v) == 0 {
		return NullUint64s{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullUint64s{}
	}

	return NewNullUint64s(string(s))
}

func (t NullUint64s) Uint64s() []uint64 {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []any
	json.Unmarshal([]byte(t.String), &v)
	if len(v) == 0 {
		return nil
	}
	w := make([]uint64, len(v))
	newV := New()
	defer newV.Close()
	for i, x := range v {
		w[i] = newV.Reload(x).DefaultUint64(0)
	}
	return w
}
func NewNullStrings(s string) NullStrings {
	var x NullStrings
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullStrings(v []string) NullStrings {
	if len(v) == 0 {
		return NullStrings{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullStrings{}
	}

	return NewNullStrings(string(s))
}
func (t NullStrings) Strings() []string {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []string
	json.Unmarshal([]byte(t.String), &v)
	return v
}
func NewNullStringMap(s string) NullStringMap {
	var x NullStringMap
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullStringMap(v map[string]string) NullStringMap {
	if len(v) == 0 {
		return NullStringMap{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullStringMap{}
	}

	return NewNullStringMap(string(s))
}

func (t NullStringMap) StringMap() map[string]string {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v map[string]string
	json.Unmarshal([]byte(t.String), &v)
	return v
}
func NewNullComplexStringMap(s string) NullComplexStringMap {
	var x NullComplexStringMap
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullComplexStringMap(v map[string]map[string]string) NullComplexStringMap {
	if len(v) == 0 {
		return NullComplexStringMap{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullComplexStringMap{}
	}

	return NewNullComplexStringMap(string(s))
}
func (t NullComplexStringMap) TStringMap() map[string]map[string]string {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v map[string]map[string]string
	json.Unmarshal([]byte(t.String), &v)
	return v
}
func NewNullStringMaps(s string) NullStringMaps {
	var x NullStringMaps
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullStringMaps(v []map[string]string) NullStringMaps {
	if len(v) == 0 {
		return NullStringMaps{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullStringMaps{}
	}

	return NewNullStringMaps(string(s))
}

func (t NullStringMaps) StringMaps() []map[string]string {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []map[string]string
	json.Unmarshal([]byte(t.String), &v)
	return v
}
func NewNullStringMapsMap(s string) NullStringMapsMap {
	var x NullStringMapsMap
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullStringMapsMap(v map[string][]map[string]string) NullStringMapsMap {
	if len(v) == 0 {
		return NullStringMapsMap{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullStringMapsMap{}
	}

	return NewNullStringMapsMap(string(s))
}
func (t NullStringMapsMap) StringMapsMap() map[string][]map[string]string {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v map[string][]map[string]string
	json.Unmarshal([]byte(t.String), &v)
	return v
}
func NewNullStringsMap(s string) NullStringsMap {
	var x NullStringsMap
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullStringsMap(v map[string][]string) NullStringsMap {
	if len(v) == 0 {
		return NullStringsMap{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullStringsMap{}
	}

	return NewNullStringsMap(string(s))
}

func (t NullStringsMap) StringsMap() map[string][]string {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v map[string][]string
	json.Unmarshal([]byte(t.String), &v)
	return v
}
func NewNullComplexStringsMap(s string) NullComplexStringsMap {
	var x NullComplexStringsMap
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullComplexStringsMap(v map[string][][]string) NullComplexStringsMap {
	if len(v) == 0 {
		return NullComplexStringsMap{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullComplexStringsMap{}
	}

	return NewNullComplexStringsMap(string(s))
}

func (t NullComplexStringsMap) StringsMap() map[string][][]string {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v map[string][][]string
	json.Unmarshal([]byte(t.String), &v)
	return v
}

func NewComplexStringMaps(s string) NullComplexStringMaps {
	var x NullComplexStringMaps
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToComplexStringMaps(v []map[string][]map[string]string) NullComplexStringMaps {
	if len(v) == 0 {
		return NullComplexStringMaps{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return NullComplexStringMaps{}
	}

	return NewComplexStringMaps(string(s))
}

func (t NullComplexStringMaps) StringMaps() []map[string][]map[string]string {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []map[string][]map[string]string
	json.Unmarshal([]byte(t.String), &v)
	return v
}

func NewNullComplexMaps(s string) ComplexMaps {
	var x ComplexMaps
	if s != "" {
		x.Scan(s)
	}
	return x
}
func ToNullComplexMaps(v []map[string]any) ComplexMaps {
	if len(v) == 0 {
		return ComplexMaps{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return ComplexMaps{}
	}

	return NewNullComplexMaps(string(s))
}

func (t ComplexMaps) Maps() []map[string]any {
	if !t.Valid || t.String == "" {
		return nil
	}
	var v []map[string]any
	json.Unmarshal([]byte(t.String), &v)
	return v
}
