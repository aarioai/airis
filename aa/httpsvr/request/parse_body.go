package request

import (
	"github.com/aarioai/airis/aa/ae"
	"mime/multipart"
	"strings"
)

func (r *Request) Files(p string, required ...bool) ([]*multipart.FileHeader, *ae.Error) {
	if !r.bodyParsed {
		e := r.parseBodyStream()
		if e != nil {
			return nil, e
		}
	}

	isRequiredBool := isRequired(required)
	if len(r.injectedFiles) == 0 {
		return nil, errorOnEmpty(p, isRequiredBool)
	}
	files, ok := r.injectedFiles[p]
	if !ok || len(files) == 0 {
		return nil, errorOnEmpty(p, isRequiredBool)
	}
	return files, nil
}
func (r *Request) File(p string, required ...bool) (*multipart.FileHeader, *ae.Error) {
	files, err := r.Files(p, required...)
	if err != nil {
		return nil, err
	}
	return files[0], nil
}

func (r *Request) BodyBool(p string) (bool, *ae.Error) {
	x, e := r.Body(p, false)
	if e != nil {
		return false, e
	}
	defer x.Close()
	return x.DefaultBool(false), e
}
func (r *Request) BodyBytes(p string) ([]byte, *ae.Error) {
	x, e := r.Body(p, false)
	if e != nil {
		return nil, e
	}
	return x.ReleaseBytes(), e
}
func (r *Request) BodyEnum(p string, required bool, validators []string) (string, *ae.Error) {
	x, e := r.BodyString(p, required)
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
func (r *Request) BodyEnum8(p string, required bool, validators []uint8) (uint8, *ae.Error) {
	x, e := r.BodyUint8(p, required)
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
func (r *Request) BodyEnum8i(p string, required bool, validators []int8) (int8, *ae.Error) {
	x, e := r.BodyInt8(p, required)
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

func (r *Request) BodyInt(p string, required ...bool) (int, *ae.Error) {
	v, e := parseInt64(r.Body, p, isRequired(required), 32)
	return int(v), e
}

func (r *Request) BodyInts(p string, required, allowZero bool) ([]int, *ae.Error) {
	values, e := r.parseInt64s(r.Body, p, required, allowZero, 32)
	return toInts(values, e)
}
func (r *Request) BodyInt8(p string, required ...bool) (int8, *ae.Error) {
	v, e := parseInt64(r.Body, p, isRequired(required), 8)
	return int8(v), e
}
func (r *Request) BodyInt8s(p string, required, allowZero bool) ([]int8, *ae.Error) {
	values, e := r.parseInt64s(r.Body, p, required, allowZero, 8)
	return toInt8s(values, e)
}
func (r *Request) BodyInt16(p string, required ...bool) (int16, *ae.Error) {
	v, e := parseInt64(r.Body, p, isRequired(required), 16)
	return int16(v), e
}
func (r *Request) BodyInt16s(p string, required, allowZero bool) ([]int16, *ae.Error) {
	values, e := r.parseInt64s(r.Body, p, required, allowZero, 16)
	return toInt16s(values, e)
}
func (r *Request) BodyInt32(p string, required ...bool) (int32, *ae.Error) {
	v, e := parseInt64(r.Body, p, isRequired(required), 32)
	return int32(v), e
}
func (r *Request) BodyInt32s(p string, required, allowZero bool) ([]int32, *ae.Error) {
	values, e := r.parseInt64s(r.Body, p, required, allowZero, 32)
	return toInt32s(values, e)
}
func (r *Request) BodyInt64(p string, required ...bool) (int64, *ae.Error) {
	return parseInt64(r.Body, p, isRequired(required), 64)
}
func (r *Request) BodyInt64s(p string, required, allowZero bool) ([]int64, *ae.Error) {
	return r.parseInt64s(r.Body, p, required, allowZero, 64)
}

func (r *Request) BodyPath(p string, required ...bool) (string, *ae.Error) {
	return r.BodyString(p, `^([\w-\/\.]+)$`, isRequired(required))
}
func (r *Request) BodyPaths(p string, required ...bool) ([]string, *ae.Error) {
	xx, e := r.BodyStrings(p, isRequired(required), false)
	if e != nil || len(xx) == 0 {
		return nil, e
	}
	paths := make([]string, len(xx))
	for i, x := range xx {
		// @TODO 平衡性能和准确性，调整到最合适的判断方法
		if x == "" || (strings.LastIndexByte(x, '.') < 0 || strings.IndexByte(x, ' ') > -1 || strings.IndexByte(x, '?') > -1 || strings.IndexByte(x, '=') > -1) {
			return nil, ae.NewBadParam(p)
		}
		paths[i] = x
	}
	return paths, e
}

// BodyStringFast 快速查询字符串
func (r *Request) BodyStringFast(p string) string {
	// false 是必须的，表示 required=false。默认 required = true
	v, _ := r.BodyString(p, false)
	return v
}

func (r *Request) BodyString(p string, required ...any) (string, *ae.Error) {
	// 不要再进行 len(params) 判断，这属于过度优化。这个函数应当优先传 params --> 不要强制，不然不利于使用
	// 如有该需求，应优先使用 QueryFast
	x, e := r.Body(p, required...)
	if e != nil {
		return "", e
	}
	return x.ReleaseString(), e
}

func (r *Request) BodyStrings(p string, required, allowEmptyString bool) ([]string, *ae.Error) {
	return r.parseStrings(r.Body, p, required, allowEmptyString)
}
func (r *Request) BodyUint(p string, required ...bool) (uint, *ae.Error) {
	v, e := parseUint64(r.Body, p, isRequired(required), 32)
	return uint(v), e
}

func (r *Request) BodyUints(p string, required, allowZero bool) ([]uint, *ae.Error) {
	values, e := r.parseUint64s(r.Body, p, required, allowZero, 32)
	return toUints(values, e)
}
func (r *Request) BodyUint8(p string, required ...bool) (uint8, *ae.Error) {
	v, e := parseUint64(r.Body, p, isRequired(required), 8)
	return uint8(v), e
}
func (r *Request) BodyUint8s(p string, required, allowZero bool) ([]uint8, *ae.Error) {
	values, e := r.parseUint64s(r.Body, p, required, allowZero, 8)
	return toUint8s(values, e)
}
func (r *Request) BodyUint16(p string, required ...bool) (uint16, *ae.Error) {
	v, e := parseUint64(r.Body, p, isRequired(required), 16)
	return uint16(v), e
}
func (r *Request) BodyUint16s(p string, required, allowZero bool) ([]uint16, *ae.Error) {
	values, e := r.parseUint64s(r.Body, p, required, allowZero, 16)
	return toUint16s(values, e)
}
func (r *Request) BodyUint32(p string, required ...bool) (uint32, *ae.Error) {
	v, e := parseUint64(r.Body, p, isRequired(required), 32)
	return uint32(v), e
}

func (r *Request) BodyUint32s(p string, required, allowZero bool) ([]uint32, *ae.Error) {
	values, e := r.parseUint64s(r.Body, p, required, allowZero, 32)
	return toUint32s(values, e)
}

func (r *Request) BodyUint64(p string, required ...bool) (uint64, *ae.Error) {
	return parseUint64(r.Body, p, isRequired(required), 64)
}

func (r *Request) BodyUint64s(p string, required, allowZero bool) ([]uint64, *ae.Error) {
	return r.parseUint64s(r.Body, p, required, allowZero, 64)
}

func (r *Request) BodyValid(p string, required bool, validator func(string) bool) (string, *ae.Error) {
	x, e := r.BodyString(p, required)
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
func (r *Request) BodyValid8(p string, required bool, validator func(uint8) bool) (uint8, *ae.Error) {
	x, e := r.BodyUint8(p, required)
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
