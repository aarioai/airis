package atype

import (
	"database/sql"
	"sync"
)

// Atype 提供类型安全和高效的类型转换
// @extend type T interface{Release()error}
type Atype struct {
	raw any
}

var (
	// 对象池，减少内存分配
	// sync.Pool 通常不需要手动释放对象，当创建的对象，没有引用时会自动回收
	atypePool = sync.Pool{
		New: func() interface{} {
			return new(Atype)
		},
	}
)

// New 从对象池获取一个*Atype实例
// 注意释放。即使不释放，也不影响。
func New(data ...any) *Atype {
	var v any
	if len(data) > 0 {
		v = data[0]
	} else {
		v = ""
	}
	a := atypePool.Get().(*Atype)
	a.raw = v
	return a
}

// Release 释放实例到对象池
// 即使这个atype 不是从对象池中获取的，也会放入对象池。不影响使用。
func (p *Atype) Release() error {
	p.raw = nil
	atypePool.Put(p)
	return nil
}

// Raw 返回原始数据
func (p *Atype) Raw() any {
	return p.raw
}

// Reload 重新加载数据
func (p *Atype) Reload(v any) *Atype {
	p.raw = v
	return p
}

// Get 从 map[string]any 中获取嵌套值
// p.Get("users.1.name") is short for p.Get("user", "1", "name")
// @warn p.Get("user", "1", "name") is diffirent with p.Get("user", 1, "name")
func (p *Atype) Get(keys ...any) (*Atype, error) {
	v, err := NewMap(p.raw).Get(keys[0], keys[1:]...)
	return New(v), err
}

// 基本类型转换方法
func (p *Atype) String() string {
	return String(p.raw)
}
func (p *Atype) ReleaseString() string {
	defer p.Release()
	return String(p.raw)
}
func (p *Atype) Bytes() []byte {
	return Bytes(p.raw)
}
func (p *Atype) ReleaseBytes() []byte {
	defer p.Release()
	return Bytes(p.raw)
}
func (p *Atype) Bool() (bool, error) {
	return Bool(p.raw)
}
func (p *Atype) ReleaseBool() (bool, error) {
	defer p.Release()
	return Bool(p.raw)
}
func (p *Atype) IsUnsigned() bool {
	return IsUnsigned(p.raw)
}

func (p *Atype) NullString() sql.NullString {
	return sql.NullString{String: p.String(), Valid: p.NotEmpty()}
}
func (p *Atype) ReleaseNullString() sql.NullString {
	defer p.Release()
	return p.NullString()
}
func (p *Atype) NullInt64() sql.NullInt64 {
	v, _ := p.Int64()
	return sql.NullInt64{Int64: v, Valid: p.NotEmpty()}
}
func (p *Atype) ReleaseNullInt64() sql.NullInt64 {
	defer p.Release()
	return p.NullInt64()
}
func (p *Atype) SqlFloat64() sql.NullFloat64 {
	v, _ := p.Float64()
	return sql.NullFloat64{Float64: v, Valid: p.NotEmpty()}
}
func (p *Atype) ReleaseNullFloat64() sql.NullFloat64 {
	defer p.Release()
	return p.SqlFloat64()
}
func (p *Atype) IsEmpty() bool {
	return IsEmpty(p.raw)
}
func (p *Atype) NotEmpty() bool {
	return NotEmpty(p.raw)
}

func (p *Atype) DefaultBool(defaultValue bool) bool {
	v, err := p.Bool()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultBool(defaultValue bool) bool {
	defer p.Release()
	return p.DefaultBool(defaultValue)
}
func (p *Atype) DefaultString(defaultValue string) string {
	v := p.String()
	if v == "" {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultString(defaultValue string) string {
	defer p.Release()
	return p.DefaultString(defaultValue)
}
func (p *Atype) DefaultBytes(defaultValue []byte) []byte {
	v := p.Bytes()
	if len(v) == 0 {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultBytes(defaultValue []byte) []byte {
	defer p.Release()
	return p.DefaultBytes(defaultValue)
}

func (p *Atype) Int8() (int8, error) {
	return Int8(p.raw)
}
func (p *Atype) ReleaseInt8() (int8, error) {
	defer p.Release()
	return p.Int8()
}
func (p *Atype) DefaultInt8(defaultValue int8) int8 {
	v, err := p.Int8()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) Int16() (int16, error) {
	return Int16(p.raw)
}
func (p *Atype) ReleaseInt16() (int16, error) {
	defer p.Release()
	return p.Int16()
}
func (p *Atype) DefaultInt16(defaultValue int16) int16 {
	v, err := p.Int16()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultInt16(defaultValue int16) int16 {
	defer p.Release()
	return p.DefaultInt16(defaultValue)
}

func (p *Atype) Int32() (int32, error) {
	return Int32(p.raw)
}
func (p *Atype) ReleaseInt32() (int32, error) {
	defer p.Release()
	return p.Int32()
}
func (p *Atype) DefaultInt32(defaultValue int32) int32 {
	v, err := p.Int32()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultInt32(defaultValue int32) int32 {
	defer p.Release()
	return p.DefaultInt32(defaultValue)
}

func (p *Atype) Int() (int, error) {
	return Int(p.raw)
}
func (p *Atype) ReleaseInt() (int, error) {
	defer p.Release()
	return p.Int()
}
func (p *Atype) DefaultInt(defaultValue int) int {
	v, err := p.Int()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultInt(defaultValue int) int {
	defer p.Release()
	return p.DefaultInt(defaultValue)
}

func (p *Atype) Int64() (int64, error) {
	return Int64(p.raw)
}
func (p *Atype) ReleaseInt64() (int64, error) {
	defer p.Release()
	return p.Int64()
}
func (p *Atype) DefaultInt64(defaultValue int64) int64 {
	v, err := p.Int64()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultInt64(defaultValue int64) int64 {
	defer p.Release()
	return p.DefaultInt64(defaultValue)
}

func (p *Atype) Uint8() (uint8, error) {
	return Uint8(p.raw)
}
func (p *Atype) ReleaseUint8() (uint8, error) {
	defer p.Release()
	return p.Uint8()
}
func (p *Atype) DefaultUint8(defaultValue uint8) uint8 {
	v, err := p.Uint8()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultUint8(defaultValue uint8) uint8 {
	defer p.Release()
	return p.DefaultUint8(defaultValue)
}

func (p *Atype) Uint16() (uint16, error) {
	return Uint16(p.raw)
}
func (p *Atype) ReleaseUint16() (uint16, error) {
	defer p.Release()
	return p.Uint16()
}

func (p *Atype) DefaultUint16(defaultValue uint16) uint16 {
	v, err := p.Uint16()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultUint16(defaultValue uint16) uint16 {
	defer p.Release()
	return p.DefaultUint16(defaultValue)
}

func (p *Atype) Uint24() (Uint24, error) {
	return Uint24b(p.raw)
}
func (p *Atype) ReleaseUint24() (Uint24, error) {
	defer p.Release()
	return p.Uint24()
}
func (p *Atype) DefaultUint24(defaultValue Uint24) Uint24 {
	v, err := p.Uint24()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultUint24(defaultValue Uint24) Uint24 {
	defer p.Release()
	return p.DefaultUint24(defaultValue)
}

func (p *Atype) Uint32() (uint32, error) {
	return Uint32(p.raw)
}
func (p *Atype) ReleaseUint32() (uint32, error) {
	defer p.Release()
	return p.Uint32()
}
func (p *Atype) DefaultUint32(defaultValue uint32) uint32 {
	v, err := p.Uint32()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultUint32(defaultValue uint32) uint32 {
	defer p.Release()
	return p.DefaultUint32(defaultValue)
}

func (p *Atype) Uint() (uint, error) {
	return Uint(p.raw)
}
func (p *Atype) ReleaseUint() (uint, error) {
	defer p.Release()
	return p.Uint()
}
func (p *Atype) DefaultUint(defaultValue uint) uint {
	v, err := p.Uint()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultUint(defaultValue uint) uint {
	defer p.Release()
	return p.DefaultUint(defaultValue)
}

func (p *Atype) Uint64() (uint64, error) {
	return Uint64(p.raw)
}

func (p *Atype) ReleaseUint64() (uint64, error) {
	defer p.Release()
	return p.Uint64()
}

func (p *Atype) DefaultUint64(defaultValue uint64) uint64 {
	v, err := p.Uint64()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultUint64(defaultValue uint64) uint64 {
	defer p.Release()
	return p.DefaultUint64(defaultValue)
}

func (p *Atype) Float32() (float32, error) {
	return Float32(p.raw)
}
func (p *Atype) ReleaseFloat32() (float32, error) {
	defer p.Release()
	return p.Float32()
}
func (p *Atype) DefaultFloat32(defaultValue float32) float32 {
	v, err := p.Float32()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultFloat32(defaultValue float32) float32 {
	defer p.Release()
	return p.DefaultFloat32(defaultValue)
}

func (p *Atype) Float64() (float64, error) {
	return Float64(p.raw, 64)
}
func (p *Atype) ReleaseFloat64() (float64, error) {
	defer p.Release()
	return p.Float64()
}
func (p *Atype) DefaultFloat64(defaultValue float64) float64 {
	v, err := p.Float64()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) ReleaseDefaultFloat64(defaultValue float64) float64 {
	defer p.Release()
	return p.DefaultFloat64(defaultValue)
}
