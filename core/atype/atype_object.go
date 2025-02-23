package atype

import (
	"encoding/json"
	"reflect"
)

func (p *Atype) IsNil() bool {
	if p.raw == nil {
		return true
	}
	v := reflect.ValueOf(p.raw)
	switch v.Kind() {
	case reflect.Ptr, reflect.UnsafePointer,
		reflect.Map, reflect.Slice,
		reflect.Array, reflect.Interface:
		return v.IsNil()
	}
	return false
}

func (p *Atype) Strings() ([]string, bool) {
	switch v := p.raw.(type) {
	case []string:
		return v, true

	case []any:
		result := make([]string, len(v))
		for i, item := range v {
			switch val := item.(type) {
			case string:
				result[i] = val
			case []byte:
				result[i] = string(val)
			default:
				return nil, false
			}
		}
		return result, true

	case [][]byte:
		result := make([]string, len(v))
		for i, bytes := range v {
			result[i] = string(bytes)
		}
		return result, true
	}

	return nil, false
}

func (p *Atype) ReleaseStrings() ([]string, bool) {
	defer p.Close()
	return p.Strings()
}

func (p *Atype) Ints() ([]int, bool) {
	switch v := p.raw.(type) {
	case []int:
		return v, true

	case []any:
		result := make([]int, len(v))
		newV := New()
		defer newV.Close()
		for i, item := range v {
			newV.Reload(item)
			val, err := newV.Int()
			if err != nil {
				return nil, false
			}
			result[i] = val
		}
		return result, true
	}

	return nil, false
}
func (p *Atype) ReleaseInts() ([]int, bool) {
	defer p.Close()
	return p.Ints()
}
func (p *Atype) Uints() ([]uint, bool) {
	switch v := p.raw.(type) {
	case []uint:
		return v, true

	case []any:
		result := make([]uint, len(v))
		newV := New()
		defer newV.Close()
		for i, item := range v {
			val, err := newV.Reload(item).Uint()
			if err != nil {
				return nil, false
			}
			result[i] = val
		}
		return result, true
	}

	return nil, false
}
func (p *Atype) ReleaseUints() ([]uint, bool) {
	defer p.Close()
	return p.Uints()
}

func (p *Atype) Int64s() ([]int64, bool) {
	switch v := p.raw.(type) {
	case []int64:
		return v, true

	case []any:
		result := make([]int64, len(v))
		newV := New()
		defer newV.Close()
		for i, item := range v {
			val, err := newV.Reload(item).Int64()
			if err != nil {
				return nil, false
			}
			result[i] = val
		}
		return result, true
	}

	return nil, false
}
func (p *Atype) ReleaseInt64s() ([]int64, bool) {
	defer p.Close()
	return p.Int64s()
}

func (p *Atype) Uint64s() ([]uint64, bool) {
	switch v := p.raw.(type) {
	case []uint64:
		return v, true

	case []any:
		result := make([]uint64, len(v))
		newV := New()
		defer newV.Close()
		for i, item := range v {
			val, err := newV.Reload(item).Uint64()
			if err != nil {
				return nil, false
			}
			result[i] = val
		}
		return result, true
	}

	return nil, false
}
func (p *Atype) ReleaseUint64s() ([]uint64, bool) {
	defer p.Close()
	return p.Uint64s()
}

func (p *Atype) Float32s() ([]float32, bool) {
	switch v := p.raw.(type) {
	case []float32:
		return v, true

	case []any:
		result := make([]float32, len(v))
		newV := New()
		defer newV.Close()
		for i, item := range v {
			val, err := newV.Reload(item).Float32()
			if err != nil {
				return nil, false
			}
			result[i] = val
		}
		return result, true
	}

	return nil, false
}
func (p *Atype) ReleaseFloat32s() ([]float32, bool) {
	defer p.Close()
	return p.Float32s()
}
func (p *Atype) Float64s() ([]float64, bool) {
	switch v := p.raw.(type) {
	case []float64:
		return v, true

	case []any:
		result := make([]float64, len(v))
		newV := New()
		defer newV.Close()
		for i, item := range v {
			val, err := newV.Reload(item).Float64()
			if err != nil {
				return nil, false
			}
			result[i] = val
		}
		return result, true
	}

	return nil, false
}
func (p *Atype) ReleaseFloat64s() ([]float64, bool) {
	defer p.Close()
	return p.Float64s()
}

func (p *Atype) ArrayJson(allowNil bool) (json.RawMessage, bool) {
	switch v := p.raw.(type) {
	case json.RawMessage:
		return v, true
	case []uint8:
		if bytes, ok := MarshalUint8s(v); ok {
			return bytes, true
		}
		return nil, false

	case []any:
		if bytes, err := json.Marshal(v); err == nil {
			return bytes, true
		}
		return nil, false
	}

	if allowNil {
		if p.IsNil() {
			return nil, true
		}
		if s, ok := p.raw.(string); ok && s == "" {
			return nil, true
		}
	}

	return nil, false
}
func (p *Atype) ReleaseArrayJson(allowNil bool) (json.RawMessage, bool) {
	defer p.Close()
	return p.ArrayJson(allowNil)
}

func (p *Atype) MapJson(allowNil bool) (json.RawMessage, bool) {
	switch v := p.raw.(type) {
	case map[string]any:
		if bytes, err := json.Marshal(v); err == nil {
			return bytes, true
		}
		return nil, false

	case []byte:
		if len(v) > 0 && v[0] == '{' {
			return v, true
		}
		return v, false
	}

	if allowNil {
		if p.IsNil() {
			return nil, true
		}
		if s, ok := p.raw.(string); ok && s == "" {
			return nil, true
		}
	}

	return nil, false
}
func (p *Atype) ReleaseMapJson(allowNil bool) (json.RawMessage, bool) {
	defer p.Close()
	return p.MapJson(allowNil)
}
