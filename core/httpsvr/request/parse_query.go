package request

import (
	"github.com/aarioai/airis/core/ae"
	"strconv"
)

func (r *Request) QueryBool(p string) (bool, *ae.Error) {
	x, e := r.Query(p, false)
	if e != nil {
		return false, e
	}
	return x.ReleaseDefaultBool(false), nil
}

func (r *Request) QueryBytes(p string) ([]byte, *ae.Error) {
	x, e := r.Query(p, false)
	if e != nil {
		return nil, e
	}
	return x.ReleaseBytes(), nil
}
func (r *Request) QueryEnum(p string, required bool, validators []string) (string, *ae.Error) {
	x, e := r.QueryString(p, required)
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
func (r *Request) QueryEnum8(p string, required bool, validators []uint8) (uint8, *ae.Error) {
	x, e := r.QueryUint8(p, required)
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
func (r *Request) QueryEnum8i(p string, required bool, validators []int8) (int8, *ae.Error) {
	x, e := r.QueryInt8(p, required)
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
func (r *Request) QueryInt(p string, required ...bool) (int, *ae.Error) {
	v, e := parseInt64(r.Query, p, isRequired(required), 32)
	return int(v), e
}
func (r *Request) QueryInts(p string, required, allowZero bool) ([]int, *ae.Error) {
	values, e := r.parseInt64s(r.Query, p, required, allowZero, 32)
	return toInts(values, e)
}
func (r *Request) QueryInt8(p string, required ...bool) (int8, *ae.Error) {
	v, e := parseInt64(r.Query, p, isRequired(required), 8)
	return int8(v), e
}
func (r *Request) QueryInt8s(p string, required, allowZero bool) ([]int8, *ae.Error) {
	values, e := r.parseInt64s(r.Query, p, required, allowZero, 8)
	return toInt8s(values, e)
}
func (r *Request) QueryInt16(p string, required ...bool) (int16, *ae.Error) {
	v, e := parseInt64(r.Query, p, isRequired(required), 16)
	return int16(v), e
}
func (r *Request) QueryInt16s(p string, required, allowZero bool) ([]int16, *ae.Error) {
	values, e := r.parseInt64s(r.Query, p, required, allowZero, 16)
	return toInt16s(values, e)
}
func (r *Request) QueryInt64(p string, required ...bool) (int64, *ae.Error) {
	return parseInt64(r.Query, p, isRequired(required), 64)
}
func (r *Request) QueryInt64s(p string, required, allowZero bool) ([]int64, *ae.Error) {
	return r.parseInt64s(r.Query, p, required, allowZero, 64)
}

// {id:uint64}  or {sid:string}
func (r *Request) QueryId(p string, params ...any) (sid string, id uint64, e *ae.Error) {
	sid, e = r.QueryString(p, params...)
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
	id, _ = strconv.ParseUint(sid, 10, 64)
	return
}

// QueryFast 更高效地快速查询字符串
func (r *Request) QueryFast(p string) string {
	// false 是必须的，表示 required=false。默认 required = true
	v, _ := r.queryString(p, false)
	return v
}

func (r *Request) QueryString(p string, params ...any) (string, *ae.Error) {
	// 不要再进行 len(params) 判断，这属于过度优化。这个函数应当优先传 params --> 不要强制，不然不利于使用
	// 如有该需求，应优先使用 QueryFast
	return r.queryString(p, params...)
}
func (r *Request) QueryStrings(p string, required, allowEmptyString bool) ([]string, *ae.Error) {
	return r.parseStrings(r.Query, p, required, allowEmptyString)
}
func (r *Request) QueryUint(p string, required ...bool) (uint, *ae.Error) {
	v, e := parseUint64(r.Query, p, isRequired(required), 32)
	return uint(v), e
}
func (r *Request) QueryUints(p string, required, allowZero bool) ([]uint, *ae.Error) {
	values, e := r.parseUint64s(r.Query, p, required, allowZero, 32)
	return toUints(values, e)
}
func (r *Request) QueryUint8(p string, required ...bool) (uint8, *ae.Error) {
	v, e := parseUint64(r.Query, p, isRequired(required), 8)
	return uint8(v), e
}
func (r *Request) QueryUint8s(p string, required, allowZero bool) ([]uint8, *ae.Error) {
	values, e := r.parseUint64s(r.Query, p, required, allowZero, 8)
	return toUint8s(values, e)
}
func (r *Request) QueryUint16(p string, required ...bool) (uint16, *ae.Error) {
	v, e := parseUint64(r.Query, p, isRequired(required), 16)
	return uint16(v), e
}
func (r *Request) QueryUint16s(p string, required, allowZero bool) ([]uint16, *ae.Error) {
	values, e := r.parseUint64s(r.Query, p, required, allowZero, 16)
	return toUint16s(values, e)
}
func (r *Request) QueryUint32(p string, required ...bool) (uint32, *ae.Error) {
	v, e := parseUint64(r.Query, p, isRequired(required), 32)
	return uint32(v), e
}
func (r *Request) QueryUint32s(p string, required, allowZero bool) ([]uint32, *ae.Error) {
	values, e := r.parseUint64s(r.Query, p, required, allowZero, 32)
	return toUint32s(values, e)
}
func (r *Request) QueryUint64(p string, required ...bool) (uint64, *ae.Error) {
	return parseUint64(r.Query, p, isRequired(required), 64)
}
func (r *Request) QueryUint64s(p string, required, allowZero bool) ([]uint64, *ae.Error) {
	return r.parseUint64s(r.Query, p, required, allowZero, 64)
}
func (r *Request) QueryValid(p string, required bool, validator func(string) bool) (string, *ae.Error) {
	x, e := r.QueryString(p, required)
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
func (r *Request) QueryValid8(p string, required bool, validator func(uint8) bool) (uint8, *ae.Error) {
	x, e := r.QueryUint8(p, required)
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
