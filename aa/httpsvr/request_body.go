package httpsvr

import (
	"github.com/aarioai/airis/aa/ae"
	"mime/multipart"
	"strings"
)

func (c *Context) Files(p string, required ...bool) ([]*multipart.FileHeader, *ae.Error) {
	if !c.bodyParsed {
		e := c.parseBodyStream()
		if e != nil {
			return nil, e
		}
	}

	isRequiredBool := isRequired(required)
	if len(c.injectedFiles) == 0 {
		return nil, errorOnEmpty(p, isRequiredBool)
	}
	files, ok := c.injectedFiles[p]
	if !ok || len(files) == 0 {
		return nil, errorOnEmpty(p, isRequiredBool)
	}
	return files, nil
}
func (c *Context) File(p string, required ...bool) (*multipart.FileHeader, *ae.Error) {
	files, err := c.Files(p, required...)
	if err != nil {
		return nil, err
	}
	return files[0], nil
}

func (c *Context) BodyBool(p string) (bool, *ae.Error) {
	x, e := c.Body(p, false)
	if e != nil {
		return false, e
	}
	defer x.Close()
	return x.DefaultBool(false), e
}
func (c *Context) BodyBytes(p string) ([]byte, *ae.Error) {
	x, e := c.Body(p, false)
	if e != nil {
		return nil, e
	}
	return x.ReleaseBytes(), e
}
func (c *Context) BodyEnum(p string, required bool, validators []string) (string, *ae.Error) {
	x, e := c.BodyString(p, required)
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
func (c *Context) BodyEnum8(p string, required bool, validators []uint8) (uint8, *ae.Error) {
	x, e := c.BodyUint8(p, required)
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
func (c *Context) BodyEnum8i(p string, required bool, validators []int8) (int8, *ae.Error) {
	x, e := c.BodyInt8(p, required)
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

func (c *Context) BodyInt(p string, required ...bool) (int, *ae.Error) {
	v, e := parseInt64(c.Body, p, isRequired(required), 32)
	return int(v), e
}

func (c *Context) BodyInts(p string, required, allowZero bool) ([]int, *ae.Error) {
	values, e := c.parseInt64s(c.Body, p, required, allowZero, 32)
	return toInts(values, e)
}
func (c *Context) BodyInt8(p string, required ...bool) (int8, *ae.Error) {
	v, e := parseInt64(c.Body, p, isRequired(required), 8)
	return int8(v), e
}
func (c *Context) BodyInt8s(p string, required, allowZero bool) ([]int8, *ae.Error) {
	values, e := c.parseInt64s(c.Body, p, required, allowZero, 8)
	return toInt8s(values, e)
}
func (c *Context) BodyInt16(p string, required ...bool) (int16, *ae.Error) {
	v, e := parseInt64(c.Body, p, isRequired(required), 16)
	return int16(v), e
}
func (c *Context) BodyInt16s(p string, required, allowZero bool) ([]int16, *ae.Error) {
	values, e := c.parseInt64s(c.Body, p, required, allowZero, 16)
	return toInt16s(values, e)
}
func (c *Context) BodyInt32(p string, required ...bool) (int32, *ae.Error) {
	v, e := parseInt64(c.Body, p, isRequired(required), 32)
	return int32(v), e
}
func (c *Context) BodyInt32s(p string, required, allowZero bool) ([]int32, *ae.Error) {
	values, e := c.parseInt64s(c.Body, p, required, allowZero, 32)
	return toInt32s(values, e)
}
func (c *Context) BodyInt64(p string, required ...bool) (int64, *ae.Error) {
	return parseInt64(c.Body, p, isRequired(required), 64)
}
func (c *Context) BodyInt64s(p string, required, allowZero bool) ([]int64, *ae.Error) {
	return c.parseInt64s(c.Body, p, required, allowZero, 64)
}

func (c *Context) BodyPath(p string, required ...bool) (string, *ae.Error) {
	return c.BodyString(p, `^([\w-\/\.]+)$`, isRequired(required))
}
func (c *Context) BodyPaths(p string, required ...bool) ([]string, *ae.Error) {
	xx, e := c.BodyStrings(p, isRequired(required), false)
	if e != nil || len(xx) == 0 {
		return nil, e
	}
	paths := make([]string, len(xx))
	for i, x := range xx {
		// @TODO Balance performance and accuracy; settle on the most suitable check.
		if x == "" || (strings.LastIndexByte(x, '.') < 0 || strings.IndexByte(x, ' ') > -1 || strings.IndexByte(x, '?') > -1 || strings.IndexByte(x, '=') > -1) {
			return nil, ae.NewBadParam(p)
		}
		paths[i] = x
	}
	return paths, e
}

// BodyStringFast looks up a string quickly.
func (c *Context) BodyStringFast(p string) string {
	// false is required: it means required=false. The default required is true.
	v, _ := c.BodyString(p, false)
	return v
}

func (c *Context) BodyString(p string, patterns ...any) (string, *ae.Error) {
	// Do not add another len(params) check; that is over-optimization. This function should accept params first — do not require them, or it becomes harder to use.
	// If that fast path is needed, prefer QueryFast.
	x, e := c.Body(p, patterns...)
	if e != nil {
		return "", e
	}
	return x.ReleaseString(), e
}

func (c *Context) BodyStrings(p string, required, allowEmptyString bool) ([]string, *ae.Error) {
	return c.parseStrings(c.Body, p, required, allowEmptyString)
}
func (c *Context) BodyUint(p string, required ...bool) (uint, *ae.Error) {
	v, e := parseUint64(c.Body, p, isRequired(required), 32)
	return uint(v), e
}

func (c *Context) BodyUints(p string, required, allowZero bool) ([]uint, *ae.Error) {
	values, e := c.parseUint64s(c.Body, p, required, allowZero, 32)
	return toUints(values, e)
}
func (c *Context) BodyUint8(p string, required ...bool) (uint8, *ae.Error) {
	v, e := parseUint64(c.Body, p, isRequired(required), 8)
	return uint8(v), e
}
func (c *Context) BodyUint8s(p string, required, allowZero bool) ([]uint8, *ae.Error) {
	values, e := c.parseUint64s(c.Body, p, required, allowZero, 8)
	return toUint8s(values, e)
}
func (c *Context) BodyUint16(p string, required ...bool) (uint16, *ae.Error) {
	v, e := parseUint64(c.Body, p, isRequired(required), 16)
	return uint16(v), e
}
func (c *Context) BodyUint16s(p string, required, allowZero bool) ([]uint16, *ae.Error) {
	values, e := c.parseUint64s(c.Body, p, required, allowZero, 16)
	return toUint16s(values, e)
}
func (c *Context) BodyUint32(p string, required ...bool) (uint32, *ae.Error) {
	v, e := parseUint64(c.Body, p, isRequired(required), 32)
	return uint32(v), e
}

func (c *Context) BodyUint32s(p string, required, allowZero bool) ([]uint32, *ae.Error) {
	values, e := c.parseUint64s(c.Body, p, required, allowZero, 32)
	return toUint32s(values, e)
}

func (c *Context) BodyUint64(p string, required ...bool) (uint64, *ae.Error) {
	return parseUint64(c.Body, p, isRequired(required), 64)
}

func (c *Context) BodyUint64s(p string, required, allowZero bool) ([]uint64, *ae.Error) {
	return c.parseUint64s(c.Body, p, required, allowZero, 64)
}

func (c *Context) BodyValid(p string, required bool, validator func(string) bool) (string, *ae.Error) {
	x, e := c.BodyString(p, required)
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
func (c *Context) BodyValid8(p string, required bool, validator func(uint8) bool) (uint8, *ae.Error) {
	x, e := c.BodyUint8(p, required)
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
