package response

import (
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/httpsvr/request"
	"github.com/aarioai/airis/pkg/types"
	"reflect"
	"strings"
)

func (w *Writer) decorateData(payload any) (any, *ae.Error) {
	tagname := w.SerializeTag
	if tagname == "" {
		tagname = SerializeTag
	}
	pf, e := w.filterFields(payload, tagname)
	if e != nil {
		return nil, e
	}

	p, e := w.stringifyBigint(pf, tagname)
	if e != nil {
		return nil, e
	}

	return p, nil
}

// @TODO
// ?_field=time,service,connections:[name,scheme],server_id,test:{a,b,c}
func (w *Writer) filterFields(a any, tagname string) (any, *ae.Error) {
	m := w.request.QueryFast(request.ParamField)
	if m == "" {
		return a, nil
	}
	if m[0] == '[' && m[len(m)-1] == ']' {
		return filterArrayFields(a, tagname, strings.Split(m[1:len(m)-1], ",")...)
	}
	return filterMapFields(a, tagname, strings.Split(m, ",")...)
}

func filterMapFields(u any, tagname string, tags ...string) (map[string]any, *ae.Error) {
	var found bool
	ret := make(map[string]any, 0)
	t := reflect.TypeOf(u)
	if t.Kind() == reflect.Map {
		for _, tag := range tags {
			found = false
			iter := reflect.ValueOf(u).MapRange()
			for iter.Next() {
				if iter.Key().String() == tag {
					found = true
					ret[tag] = iter.Value().Interface()
				}
			}
			if !found {
				ret[tag] = nil
			}
		}
		return ret, nil
	}

	for _, tag := range tags {
		found = false
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			al := f.Tag.Get(tagname)
			if al == tag {
				found = true
				ret[tag] = reflect.ValueOf(u).FieldByName(f.Name).Interface()
			}
		}
		if !found {
			ret[tag] = nil
		}
	}
	return ret, nil
}

func filterArrayFields(w any, tagname string, tags ...string) (ret []map[string]any, e *ae.Error) {
	t := reflect.TypeOf(w).Kind()
	if t != reflect.Slice && t != reflect.Array {
		return nil, ae.NewBadParam(request.ParamField)
	}
	v := reflect.ValueOf(w)
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
func (w *Writer) stringifyBigint(payload any, tagname string) (any, *ae.Error) {
	stringify, _ := w.request.QueryBool(request.ParamStringify)
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
	// v 有可能是一个nil指针
	if t.NumField() == 0 {
		return nil, nil
	}

	p := make(map[string]any, 0)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		ks := f.Tag.Get(tagname)
		// 忽略json/xml 里面的  -
		if ks == "-" {
			continue
		}
		v1 := v.FieldByName(f.Name)
		if !v1.CanInterface() {
			continue
		}
		w, e := StringifyBigintFields(v1.Interface(), tagname)
		if e != nil {
			return nil, e
		}

		//   struct in struct
		if ks == "" {
			m, ok := w.(map[string]any)
			if !ok {
				//p[t.Name()] = w  // 忽略 json/xml tag不存在或为空
				continue
			} else {
				for y, z := range m {
					p[y], _ = StringifyBigintFields(z, tagname)
				}
			}

		} else {
			p[ks] = w
		}

	}
	return p, nil
}
func stringifyMap(t reflect.Type, v reflect.Value, tagname string) (any, *ae.Error) {
	// v 有可能是一个nil指针
	if len(v.MapKeys()) == 0 {
		return nil, nil
	}
	p := make(map[string]any, v.Len())
	for _, key := range v.MapKeys() {
		ks := key.String()
		// 忽略json/xml 里面的  -
		if ks == "-" {
			continue
		}
		w, e := StringifyBigintFields(v.MapIndex(key).Interface(), tagname)
		if e != nil {
			return nil, e
		}

		//   struct in struct
		if ks == "" {
			m, ok := w.(map[string]any)
			if !ok {
				p[t.Name()] = w
			} else {
				for y, z := range m {
					p[y], _ = StringifyBigintFields(z, tagname)
				}
			}
		} else {
			p[ks] = w
		}
	}
	return p, nil
}

// 2.0 版本，仅针对初级原始类型为 int64/uint64 字段转为 string
func StringifyBigintFields(payload any, tagname string) (any, *ae.Error) {
	if payload == nil {
		return nil, nil
	}

	t := reflect.TypeOf(payload)
	v := reflect.ValueOf(payload)
	k := v.Kind()
	// 指针
	if k == reflect.Ptr {
		v = v.Elem() // 必须放t=v.Type前面
		t = t.Elem() // 必须用 t.Elem()，不能用 v.Type()
		k = v.Kind()
	}
	if k == reflect.Interface {
		k = v.Kind()
		if k == reflect.Ptr {
			v = v.Elem() // 必须放t=v.Type前面
			t = t.Elem() // 必须用 t.Elem()，不能用 v.Type()
			k = v.Kind()
		}
	}
	if k == reflect.Invalid {
		return nil, nil
	}
	switch k {
	case reflect.Invalid:
		return nil, nil // v 有可能是一个nil指针
	case reflect.Slice, reflect.Array:
		return stringifySlice(v, tagname)
	case reflect.Struct:
		return stringifyStruct(t, v, tagname)
	case reflect.Map:
		return stringifyMap(t, v, tagname)
	case reflect.Int64: // 也有可能是 int64 变体，所以不能用 x.(int64) 转换，应该用强制类型转换
		return types.FormatInt(v.Int()), nil
	case reflect.Uint64:
		return types.FormatUint(v.Uint()), nil
	}
	return payload, nil
}
