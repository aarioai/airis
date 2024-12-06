package atype

import (
	"database/sql"
	"sync"
)

// Atype 提供类型安全和高效的类型转换
type Atype struct {
	raw any
}

var (
	// 对象池，减少内存分配
	atypePool = sync.Pool{
		New: func() interface{} {
			return new(Atype)
		},
	}
)

// New 创建一个新的 Atype 实例
// 即使不释放，也不影响。
func New(data any) *Atype {
	a := atypePool.Get().(*Atype)
	a.raw = data
	return a
}

// Release 释放 Atype 实例到对象池
// 即使这个atype 不是从对象池中获取的，也会放入对象池。不影响使用。
func (p *Atype) Release() {
	p.raw = nil
	atypePool.Put(p)
}

// Raw 返回原始数据
func (p *Atype) Raw() any {
	return p.raw
}

// Reload 重新加载数据
func (p *Atype) Reload(v any) {
	p.raw = v
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
func (p *Atype) Bytes() []byte {
	return Bytes(p.raw)
}
func (p *Atype) Slice() ([]any, error) {
	return Slice(p.raw)
}
func (p *Atype) Bool() (bool, error) {
	return Bool(p.raw)
}

func (p *Atype) SqlNullString() sql.NullString {
	return sql.NullString{String: p.String(), Valid: p.NotEmpty()}
}

func (p *Atype) SqlNullInt64() sql.NullInt64 {
	v, _ := p.Int64()
	return sql.NullInt64{Int64: v, Valid: p.NotEmpty()}
}

func (p *Atype) SqlNullFloat64() sql.NullFloat64 {
	v, _ := p.Float64()
	return sql.NullFloat64{Float64: v, Valid: p.NotEmpty()}
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

func (p *Atype) DefaultSlice(defaultValue []any) []any {
	v, err := p.Slice()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) DefaultString(defaultValue string) string {
	v := p.String()
	if v == "" {
		return defaultValue
	}
	return v
}

func (p *Atype) DefaultBytes(defaultValue []byte) []byte {
	v := p.Bytes()
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func (p *Atype) Int8() (int8, error) {
	return Int8(p.raw)
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
func (p *Atype) DefaultInt16(defaultValue int16) int16 {
	v, err := p.Int16()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) Int32() (int32, error) {
	return Int32(p.raw)
}

func (p *Atype) DefaultInt32(defaultValue int32) int32 {
	v, err := p.Int32()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) Int() (int, error) {
	return Int(p.raw)
}

func (p *Atype) DefaultInt(defaultValue int) int {
	v, err := p.Int()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) Int64() (int64, error) {
	return Int64(p.raw)
}

func (p *Atype) DefaultInt64(defaultValue int64) int64 {
	v, err := p.Int64()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) Uint8() (uint8, error) {
	return Uint8(p.raw)
}

func (p *Atype) DefaultUint8(defaultValue uint8) uint8 {
	v, err := p.Uint8()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) Uint16() (uint16, error) {
	return Uint16(p.raw)
}

func (p *Atype) DefaultUint16(defaultValue uint16) uint16 {
	v, err := p.Uint16()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) Uint24() (Uint24, error) {
	return Uint24b(p.raw)
}

func (p *Atype) DefaultUint24(defaultValue Uint24) Uint24 {
	v, err := p.Uint24()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) Uint32() (uint32, error) {
	return Uint32(p.raw)
}

func (p *Atype) DefaultUint32(defaultValue uint32) uint32 {
	v, err := p.Uint32()
	if err != nil {
		return defaultValue
	}
	return v
}
func (p *Atype) Uint() (uint, error) {
	return Uint(p.raw)
}

func (p *Atype) DefaultUint(defaultValue uint) uint {
	v, err := p.Uint()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) Uint64() (uint64, error) {
	return Uint64(p.raw)
}

func (p *Atype) DefaultUint64(defaultValue uint64) uint64 {
	v, err := p.Uint64()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) Float32() (float32, error) {
	return Float32(p.raw)
}

func (p *Atype) DefaultFloat32(defaultValue float32) float32 {
	v, err := p.Float32()
	if err != nil {
		return defaultValue
	}
	return v
}

func (p *Atype) Float64() (float64, error) {
	return Float64(p.raw, 64)
}

func (p *Atype) DefaultFloat64(defaultValue float64) float64 {
	v, err := p.Float64()
	if err != nil {
		return defaultValue
	}
	return v
}
