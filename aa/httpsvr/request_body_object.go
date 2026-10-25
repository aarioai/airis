package httpsvr

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
)

func (c *Context) BodyAnyMap(p string, requireds ...bool) (map[string]any, *ae.Error) {
	required := len(requireds) == 0 || requireds[0]
	x, e := c.Body(p, required)
	if e != nil {
		return nil, e
	}
	defer x.Close()
	if x.IsNil() || x.String() == "" {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	raw := x.Raw()
	b, ok := raw.(map[string]any)
	if !ok {
		return nil, ae.NewBadParam(p)
	}
	if len(b) == 0 {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	return b, nil
}
func (c *Context) BodyFloat64Map(p string, requireds ...bool) (map[string]float64, *ae.Error) {
	b, e := c.BodyAnyMap(p, requireds...)
	if e != nil {
		return nil, e
	}
	maps, err := atype.ToFloat64Map(b)
	if err != nil {
		return nil, ae.NewBadParam(p)
	}
	required := len(requireds) == 0 || requireds[0]
	if required && maps == nil {
		return nil, ae.NewBadParam(p)
	}
	return maps, nil
}
func (c *Context) BodyAnySlice(p string, requireds ...bool) ([]any, *ae.Error) {
	required := len(requireds) == 0 || requireds[0]
	x, e := c.Body(p, required)
	if e != nil {
		return nil, e
	}
	defer x.Close()
	if x.IsNil() || x.String() == "" {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	raw := x.Raw()
	b, ok := raw.([]any)
	if !ok {
		return nil, ae.NewBadParam(p)
	}
	if len(b) == 0 {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	return b, nil
}

func (c *Context) BodyComplexStringMap(p string, requireds ...bool) (map[string]map[string]string, *ae.Error) {
	b, e := c.BodyAnyMap(p, requireds...)
	if e != nil {
		return nil, e
	}
	maps := atype.ToComplexStringMap(b)
	required := len(requireds) == 0 || requireds[0]
	if required && maps == nil {
		return nil, ae.NewBadParam(p)
	}
	return maps, nil
}
func (c *Context) BodyComplexStringsMap(p string, requireds ...bool) (map[string][][]string, *ae.Error) {
	b, e := c.BodyAnyMap(p, requireds...)
	if e != nil {
		return nil, e
	}
	maps := atype.ToComplexStringsMap(b)
	required := len(requireds) == 0 || requireds[0]
	if required && maps == nil {
		return nil, ae.NewBadParam(p)
	}
	return maps, nil
}
func (c *Context) BodyConvStringMaps(p string, requireds ...bool) ([]map[string]string, *ae.Error) {
	b, e := c.BodyAnySlice(p, requireds...)
	if e != nil {
		return nil, e
	}
	maps := atype.ToStringMaps(b)
	required := len(requireds) == 0 || requireds[0]
	if required && maps == nil {
		return nil, ae.NewBadParam(p)
	}
	return maps, nil
}

func (c *Context) BodyComplexMaps(p string, requireds ...bool) ([]map[string]any, *ae.Error) {
	b, e := c.BodyAnySlice(p, requireds...)
	if e != nil {
		return nil, e
	}
	maps := atype.ToComplexMaps(b)
	required := len(requireds) == 0 || requireds[0]
	if required && maps == nil {
		return nil, ae.NewBadParam(p)
	}
	return maps, nil
}
func (c *Context) BodyStringMap(p string, requireds ...bool) (map[string]string, *ae.Error) {
	b, e := c.BodyAnyMap(p, requireds...)
	if e != nil {
		return nil, e
	}
	maps := atype.ToStringMap(b)
	required := len(requireds) == 0 || requireds[0]
	if required && maps == nil {
		return nil, ae.NewBadParam(p)
	}
	return maps, nil
}
func (c *Context) BodyStringsMap(p string, requireds ...bool) (map[string][]string, *ae.Error) {
	b, e := c.BodyAnyMap(p, requireds...)
	if e != nil {
		return nil, e
	}
	maps := atype.ToStringsMap(b)
	required := len(requireds) == 0 || requireds[0]
	if required && maps == nil {
		return nil, ae.NewBadParam(p)
	}
	return maps, nil
}
