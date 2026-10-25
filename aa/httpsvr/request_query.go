package httpsvr

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
)

func (c *Context) QueryBool(p string) (bool, *ae.Error) {
	x, e := c.Query(p, false)
	if e != nil {
		return false, e
	}
	return x.ReleaseDefaultBool(false), nil
}

func (c *Context) QueryBytes(p string) ([]byte, *ae.Error) {
	x, e := c.Query(p, false)
	if e != nil {
		return nil, e
	}
	return x.ReleaseBytes(), nil
}
func (c *Context) QueryEnum(p string, required bool, validators []string) (string, *ae.Error) {
	x, e := c.QueryString(p, required)
	if e != nil {
		return "", e
	}
	if !required && x == "" {
		return "", nil
	}
	for _, val := range validators {
		if x == val {
			return x, nil
		}
	}
	return x, ae.NewBadParam(p)
}
func (c *Context) QueryEnum8(p string, required bool, validators []uint8) (uint8, *ae.Error) {
	x, e := c.QueryUint8(p, required)
	if e != nil {
		return 0, e
	}
	if !required && x == 0 {
		return 0, nil
	}
	for _, val := range validators {
		if x == val {
			return x, nil
		}
	}
	return x, ae.NewBadParam(p)
}
func (c *Context) QueryEnum8i(p string, required bool, validators []int8) (int8, *ae.Error) {
	x, e := c.QueryInt8(p, required)
	if e != nil {
		return 0, e
	}
	if !required && x == 0 {
		return 0, nil
	}
	for _, val := range validators {
		if x == val {
			return x, nil
		}
	}
	return x, ae.NewBadParam(p)
}
func (c *Context) QueryInt(p string, required ...bool) (int, *ae.Error) {
	v, e := parseInt64(c.Query, p, isRequired(required), 32)
	return int(v), e
}
func (c *Context) QueryInts(p string, required, allowZero bool) ([]int, *ae.Error) {
	values, e := c.parseInt64s(c.Query, p, required, allowZero, 32)
	return toInts(values, e)
}
func (c *Context) QueryInt8(p string, required ...bool) (int8, *ae.Error) {
	v, e := parseInt64(c.Query, p, isRequired(required), 8)
	return int8(v), e
}
func (c *Context) QueryInt8s(p string, required, allowZero bool) ([]int8, *ae.Error) {
	values, e := c.parseInt64s(c.Query, p, required, allowZero, 8)
	return toInt8s(values, e)
}
func (c *Context) QueryInt16(p string, required ...bool) (int16, *ae.Error) {
	v, e := parseInt64(c.Query, p, isRequired(required), 16)
	return int16(v), e
}
func (c *Context) QueryInt16s(p string, required, allowZero bool) ([]int16, *ae.Error) {
	values, e := c.parseInt64s(c.Query, p, required, allowZero, 16)
	return toInt16s(values, e)
}
func (c *Context) QueryInt64(p string, required ...bool) (int64, *ae.Error) {
	return parseInt64(c.Query, p, isRequired(required), 64)
}
func (c *Context) QueryInt64s(p string, required, allowZero bool) ([]int64, *ae.Error) {
	return c.parseInt64s(c.Query, p, required, allowZero, 64)
}

// {id:uint64}  or {sid:string}
func (c *Context) QueryId(p string, params ...any) (sid string, id uint64, e *ae.Error) {
	sid, e = c.QueryString(p, params...)
	if sid == "" || sid == "0" {
		return
	}
	for _, s := range sid {
		if s < '0' || (s > '9' && s < 'A') || (s > 'Z' && s < '_') || (s > '_' && s < 'a') || s > 'z' {
			e = ae.NewBadParam(p)
			return
		}
		if s > '9' {
			return
		}
	}
	id = types.ToUint64(sid)
	return
}

// QueryFast looks up a string more efficiently.
func (c *Context) QueryFast(p string) string {
	// false is required: it means required=false. The default required is true.
	v, _ := c.queryString(p, false)
	return v
}

func (c *Context) QueryString(p string, params ...any) (string, *ae.Error) {
	// Do not add another len(params) check; that is over-optimization. This function should accept params first — do not require them, or it becomes harder to use.
	// If that fast path is needed, prefer QueryFast.
	return c.queryString(p, params...)
}
func (c *Context) QueryStrings(p string, required, allowEmptyString bool) ([]string, *ae.Error) {
	return c.parseStrings(c.Query, p, required, allowEmptyString)
}
func (c *Context) QueryUint(p string, required ...bool) (uint, *ae.Error) {
	v, e := parseUint64(c.Query, p, isRequired(required), 32)
	return uint(v), e
}
func (c *Context) QueryUints(p string, required, allowZero bool) ([]uint, *ae.Error) {
	values, e := c.parseUint64s(c.Query, p, required, allowZero, 32)
	return toUints(values, e)
}
func (c *Context) QueryUint8(p string, required ...bool) (uint8, *ae.Error) {
	v, e := parseUint64(c.Query, p, isRequired(required), 8)
	return uint8(v), e
}
func (c *Context) QueryUint8s(p string, required, allowZero bool) ([]uint8, *ae.Error) {
	values, e := c.parseUint64s(c.Query, p, required, allowZero, 8)
	return toUint8s(values, e)
}
func (c *Context) QueryUint16(p string, required ...bool) (uint16, *ae.Error) {
	v, e := parseUint64(c.Query, p, isRequired(required), 16)
	return uint16(v), e
}
func (c *Context) QueryUint16s(p string, required, allowZero bool) ([]uint16, *ae.Error) {
	values, e := c.parseUint64s(c.Query, p, required, allowZero, 16)
	return toUint16s(values, e)
}
func (c *Context) QueryUint32(p string, required ...bool) (uint32, *ae.Error) {
	v, e := parseUint64(c.Query, p, isRequired(required), 32)
	return uint32(v), e
}
func (c *Context) QueryUint32s(p string, required, allowZero bool) ([]uint32, *ae.Error) {
	values, e := c.parseUint64s(c.Query, p, required, allowZero, 32)
	return toUint32s(values, e)
}
func (c *Context) QueryUint64(p string, required ...bool) (uint64, *ae.Error) {
	return parseUint64(c.Query, p, isRequired(required), 64)
}
func (c *Context) QueryUint64s(p string, required, allowZero bool) ([]uint64, *ae.Error) {
	return c.parseUint64s(c.Query, p, required, allowZero, 64)
}
func (c *Context) QueryValid(p string, required bool, validator func(string) bool) (string, *ae.Error) {
	x, e := c.QueryString(p, required)
	if e != nil {
		return "", e
	}
	if !required && x == "" {
		return "", nil
	}
	if ok := validator(x); !ok {
		return "", ae.NewBadParam(p)
	}
	return x, nil
}
func (c *Context) QueryValid8(p string, required bool, validator func(uint8) bool) (uint8, *ae.Error) {
	x, e := c.QueryUint8(p, required)
	if e != nil {
		return 0, e
	}
	if !required && x == 0 {
		return 0, nil
	}
	if ok := validator(x); !ok {
		return 0, ae.NewBadParam(p)
	}
	return x, nil
}
