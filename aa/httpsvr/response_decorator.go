package httpsvr

import (
	"github.com/aarioai/airis/aa/ae"
		"github.com/aarioai/airis/pkg/types"
	"reflect"
	"strings"
)

func (c *Context) decorateData(payload any) (any, *ae.Error) {
	tagname := c.SerializeTag
	if tagname == "" {
		tagname = SerializeTag
	}
	pf, e := c.filterFields(payload, tagname)
	if e != nil {
		return nil, e
	}

	p, e := c.stringifyBigint(pf, tagname)
	if e != nil {
		return nil, e
	}

	return p, nil
}

// @TODO
// ?_field=time,service,connections:[name,scheme],server_id,test:{a,b,c}
func (c *Context) filterFields(a any, tagname string) (any, *ae.Error) {
	m := c.QueryFast(ParamField)
	if m == "" {
		return a, nil
	}
	if m[0] == '[' && m[len(m)-1] == ']' {
		return filterArrayFields(a, tagname, strings.Split(m[1:len(m)-1], ",")...)
	}
	return filterMapFields(a, tagname, strings.Split(m, ",")...)
}

// structTagName extracts name from `json:"name,omitempty"`.
// inline means there is no name and the embedded struct should be expanded; skip means the tag is "-".
func structTagName(tag string) (name string, inline, skip bool) {
	if tag == "" {
		return "", true, false
	}
	name, _, _ = strings.Cut(tag, ",")
	if name == "-" {
		return "", false, true
	}
	if name == "" {
		return "", true, false
	}
	return name, false, false
}

func derefValue(u any) (reflect.Type, reflect.Value, bool) {
	if u == nil {
		return nil, reflect.Value{}, false
	}
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)
	for t.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, reflect.Value{}, false
		}
		t = t.Elem()
		v = v.Elem()
	}
	return t, v, true
}

func filterMapFields(u any, tagname string, tags ...string) (map[string]any, *ae.Error) {
	ret := make(map[string]any, len(tags))
	t, v, ok := derefValue(u)
	if !ok {
		for _, tag := range tags {
			ret[tag] = nil
		}
		return ret, nil
	}
	lookup := make(map[string]any, len(tags))
	switch t.Kind() {
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			lookup[iter.Key().String()] = iter.Value().Interface()
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			name, _, skip := structTagName(f.Tag.Get(tagname))
			if skip || name == "" {
				continue
			}
			fv := v.Field(i)
			if !fv.CanInterface() {
				continue
			}
			lookup[name] = fv.Interface()
		}
	default:
		return nil, ae.NewBadParam(ParamField)
	}
	for _, tag := range tags {
		if val, found := lookup[tag]; found {
			ret[tag] = val
		} else {
			ret[tag] = nil
		}
	}
	return ret, nil
}

func filterArrayFields(w any, tagname string, tags ...string) (ret []map[string]any, e *ae.Error) {
	_, v, ok := derefValue(w)
	if !ok {
		return nil, ae.NewBadParam(ParamField)
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, ae.NewBadParam(ParamField)
	}
	ret = make([]map[string]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		ret[i], e = filterMapFields(v.Index(i).Interface(), tagname, tags...)
		if e != nil {
			return nil, e
		}
	}

	return ret, nil
}

// ?_stringify=1  weak language, turn int64/uint64 fields into string
func (c *Context) stringifyBigint(payload any, tagname string) (any, *ae.Error) {
	var stringify bool
	s := c.HeaderFast(HeaderStringify)
	if s != "" {
		stringify = s == "1" || s == "true" || s == "TRUE" || s == "True"
	} else {
		stringify, _ = c.QueryBool(ParamStringify)
	}

	if stringify {
		return StringifyBigintFields(payload, tagname)
	}
	return payload, nil
}

func stringifySlice(v reflect.Value, tagname string) (any, *ae.Error) {
	if v.Len() == 0 {
		return nil, nil
	}
	p := make([]any, v.Len())
	var e *ae.Error
	for i := 0; i < v.Len(); i++ {
		if !v.Index(i).CanInterface() {
			return nil, nil
		}
		p[i], e = StringifyBigintFields(v.Index(i).Interface(), tagname)
		if e != nil {
			return nil, e
		}
	}
	return p, nil
}
func stringifyStruct(t reflect.Type, v reflect.Value, tagname string) (any, *ae.Error) {
	// v may be a nil pointer
	if t.NumField() == 0 {
		return nil, nil
	}

	p := make(map[string]any, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name, inline, skip := structTagName(f.Tag.Get(tagname))
		if skip {
			continue
		}
		v1 := v.Field(i)
		if !v1.CanInterface() {
			continue
		}
		converted, e := StringifyBigintFields(v1.Interface(), tagname)
		if e != nil {
			return nil, e
		}

		// Expand the embedded struct. The field itself is already converted, so do not recurse again.
		if inline {
			m, ok := converted.(map[string]any)
			if !ok {
				continue
			}
			for y, z := range m {
				p[y] = z
			}
			continue
		}
		p[name] = converted
	}
	return p, nil
}
func stringifyMap(t reflect.Type, v reflect.Value, tagname string) (any, *ae.Error) {
	if v.Len() == 0 {
		return nil, nil
	}
	p := make(map[string]any, v.Len())
	iter := v.MapRange()
	for iter.Next() {
		ks := iter.Key().String()
		// Ignore "-" keys inside json/xml.
		if ks == "-" {
			continue
		}
		converted, e := StringifyBigintFields(iter.Value().Interface(), tagname)
		if e != nil {
			return nil, e
		}

		if ks == "" {
			m, ok := converted.(map[string]any)
			if !ok {
				p[t.Name()] = converted
				continue
			}
			for y, z := range m {
				p[y] = z
			}
			continue
		}
		p[ks] = converted
	}
	return p, nil
}

// Version 2.0: convert only fields whose elementary kind is int64 or uint64 to string.
func StringifyBigintFields(payload any, tagname string) (any, *ae.Error) {
	if payload == nil {
		return nil, nil
	}

	t := reflect.TypeOf(payload)
	v := reflect.ValueOf(payload)
	k := v.Kind()
	// Unwrap pointers and interfaces. Elem of a nil pointer is Invalid and must not have Type called on it.
	for k == reflect.Pointer || k == reflect.Interface {
		if v.IsNil() {
			return nil, nil
		}
		if k == reflect.Pointer {
			t = t.Elem()
		} else {
			t = v.Elem().Type()
		}
		v = v.Elem()
		k = v.Kind()
	}
	if k == reflect.Invalid {
		return nil, nil
	}
	switch k {
	case reflect.Invalid:
		return nil, nil // v may be a nil pointer
	case reflect.Slice, reflect.Array:
		return stringifySlice(v, tagname)
	case reflect.Struct:
		return stringifyStruct(t, v, tagname)
	case reflect.Map:
		return stringifyMap(t, v, tagname)
	case reflect.Int64: // may be a defined int64 type, so do not use x.(int64); use a conversion
		return types.FormatInt(v.Int()), nil
	case reflect.Uint64:
		return types.FormatUint(v.Uint()), nil
	}
	return payload, nil
}
