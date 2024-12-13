package request

import (
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/atype"
)

func (r *Request) BodyAnyMap(p string, requireds ...bool) (map[string]any, *ae.Error) {
	required := len(requireds) == 0 || requireds[0]
	x, e := r.Body(p, required)
	if e != nil {
		return nil, e
	}
	defer x.Release()
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
func (r *Request) BodyFloat64Map(p string, requireds ...bool) (map[string]float64, *ae.Error) {
	b, e := r.BodyAnyMap(p, requireds...)
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
func (r *Request) BodyAnySlice(p string, requireds ...bool) ([]any, *ae.Error) {
	required := len(requireds) == 0 || requireds[0]
	x, e := r.Body(p, required)
	if e != nil {
		return nil, e
	}
	defer x.Release()
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

func (r *Request) BodyComplexStringMap(p string, requireds ...bool) (map[string]map[string]string, *ae.Error) {
	b, e := r.BodyAnyMap(p, requireds...)
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
func (r *Request) BodyComplexStringsMap(p string, requireds ...bool) (map[string][][]string, *ae.Error) {
	b, e := r.BodyAnyMap(p, requireds...)
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
func (r *Request) BodyConvStringMaps(p string, requireds ...bool) ([]map[string]string, *ae.Error) {
	b, e := r.BodyAnySlice(p, requireds...)
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

func (r *Request) BodyComplexMaps(p string, requireds ...bool) ([]map[string]any, *ae.Error) {
	b, e := r.BodyAnySlice(p, requireds...)
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
func (r *Request) BodyStringMap(p string, requireds ...bool) (map[string]string, *ae.Error) {
	b, e := r.BodyAnyMap(p, requireds...)
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
func (r *Request) BodyStringsMap(p string, requireds ...bool) (map[string][]string, *ae.Error) {
	b, e := r.BodyAnyMap(p, requireds...)
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
