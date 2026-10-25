package httpsvr

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"reflect"
	"regexp"
	"strings"
)

// newlineSep treats consecutive newlines as separators. Compiled once so each parameter split does not rebuild the automaton.
var newlineSep = regexp.MustCompile(`\s*[\r\n]+\s*`)

func parseInt64(method func(string, ...any) (*RawValue, *ae.Error), p string, required bool, bitSize int) (int64, *ae.Error) {
	value, e := method(p, required)
	if e != nil {
		return 0, e
	}
	s := value.String()
	if s == "" {
		if required {
			return 0, ae.NewBadParam(p)
		}
		return 0, nil
	}
	v, err := types.ParseInt64(s, bitSize)
	if err != nil {
		return 0, ae.NewBadParam(p)
	}
	return v, nil
}
func parseUint64(method func(string, ...any) (*RawValue, *ae.Error), p string, required bool, bitSize int) (uint64, *ae.Error) {
	value, e := method(p, required)
	if e != nil {
		return 0, e
	}
	s := value.String()
	if s == "" {
		if required {
			return 0, ae.NewBadParam(p)
		}
		return 0, nil
	}
	v, err := types.ParseUint64(s)
	if err != nil {
		return 0, ae.NewBadParam(p)
	}
	return v, nil
}
func (c *Context) parseInt64s(method func(string, ...any) (*RawValue, *ae.Error), p string, required, allowZero bool, bitSize int) ([]int64, *ae.Error) {
	q, e := method(p, required)
	if e != nil {
		return nil, e
	}
	var v []int64
	d := q.Raw()
	if d == nil {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	if _, ok := d.(string); ok {
		return c.sepInt64s(method, p, ",", required, allowZero, bitSize)
	}
	if reflect.TypeOf(d).Kind() != reflect.Slice {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	// May be a mixed array such as [1,"2",3].
	s := reflect.ValueOf(d)
	v = make([]int64, 0, s.Len())
	var n int64
	var err error
	for i := 0; i < s.Len(); i++ {
		n, err = atype.Int64Base(s.Index(i).Interface(), bitSize)
		if err != nil {
			return nil, ae.NewBadParam(p)
		}
		if allowZero || n > 0 {
			v = append(v, n)
		}
	}
	if len(v) == 0 && required {
		return nil, ae.NewBadParam(p)
	}
	return v, nil
}
func (c *Context) parseUint64s(method func(string, ...any) (*RawValue, *ae.Error), p string, required, allowZero bool, bitSize int) ([]uint64, *ae.Error) {
	q, e := method(p, required)
	if e != nil {
		return nil, e
	}
	var v []uint64
	d := q.Raw()
	if d == nil {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	if _, ok := d.(string); ok {
		return c.separatedUint64s(method, p, ",", required, allowZero, bitSize)
	}
	if reflect.TypeOf(d).Kind() != reflect.Slice {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	// May be a mixed array such as [1,"2",3].
	s := reflect.ValueOf(d)
	v = make([]uint64, 0, s.Len())
	var n uint64
	var err error
	for i := 0; i < s.Len(); i++ {
		n, err = atype.Uint64Base(s.Index(i).Interface(), bitSize)
		if err != nil {
			return nil, ae.NewBadParam(p)
		}
		if allowZero || n > 0 {
			v = append(v, n)
		}
	}
	if len(v) == 0 && required {
		return nil, ae.NewBadParam(p)
	}
	return v, nil
}
func (c *Context) parseStrings(method func(string, ...any) (*RawValue, *ae.Error), p string, required, allowEmptyString bool) ([]string, *ae.Error) {
	q, e := method(p, required)
	if e != nil {
		return nil, e
	}
	defer q.Close()
	var v []string
	d := q.Raw()
	if d == nil {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	if _, ok := d.(string); ok {
		return splitBySep(q.String(), ",", p, required, allowEmptyString)
	}

	if reflect.TypeOf(d).Kind() == reflect.Slice {
		// May be a mixed array such as [1,"2",3].
		s := reflect.ValueOf(d)
		v = make([]string, 0, s.Len())
		for i := 0; i < s.Len(); i++ {
			// Do not use s.Index(i).String(); it returns: <interface {} Value>
			ts := atype.String(s.Index(i).Interface())
			if allowEmptyString || (ts != "" && ts != "0") {
				v = append(v, ts)
			}
		}
	}

	if len(v) == 0 && required {
		return nil, ae.NewBadParam(p)
	}
	return v, nil
}

// @warn do not pass float
func sepStrings(method func(string, ...any) (*RawValue, *ae.Error), p string, sep string, required, allowEmptyString bool) ([]string, *ae.Error) {
	s, e := method(p, required)
	if e != nil {
		return nil, e
	}
	if s == nil {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	defer s.Close()
	return splitBySep(s.String(), sep, p, required, allowEmptyString)
}

func splitBySep(s, sep, p string, required, allowEmptyString bool) ([]string, *ae.Error) {
	// Treat newlines as separators.
	x := newlineSep.ReplaceAllString(s, sep)
	arr := strings.Split(x, sep)
	b := make([]string, 0, len(arr))
	for _, a := range arr {
		if !allowEmptyString {
			a = strings.Trim(a, " ")
			if a == "" {
				continue
			}
		}
		b = append(b, a)
	}
	if len(b) == 0 && required {
		return nil, ae.NewBadParam(p)
	}
	return b, nil
}

// @warn do not pass float
func (c *Context) sepInt64s(method func(string, ...any) (*RawValue, *ae.Error), p string, sep string, required, allowZero bool, bitSize int) ([]int64, *ae.Error) {
	arr, e := sepStrings(method, p, sep, required, false)
	if e != nil {
		return nil, e
	}
	v := make([]int64, 0, len(arr))
	var err error
	var n int64
	for _, a := range arr {
		n, err = types.ParseInt64(a, bitSize)
		if err != nil {
			return nil, ae.NewBadParam(p)
		}
		if allowZero || n > 0 {
			v = append(v, n)
		}
	}
	return v, nil
}
func (c *Context) separatedUint64s(method func(string, ...any) (*RawValue, *ae.Error), p string, sep string, required, allowZero bool, bitSize int) ([]uint64, *ae.Error) {
	arr, e := sepStrings(method, p, sep, required, false)
	if e != nil {
		return nil, e
	}
	v := make([]uint64, 0, len(arr))
	var err error
	var n uint64
	for _, a := range arr {
		n, err = types.ParseUint64(a, bitSize)
		if err != nil {
			return nil, ae.NewBadParam(p)
		}
		if allowZero || n > 0 {
			v = append(v, n)
		}
	}
	return v, nil
}
func toInts(values []int64, e *ae.Error) ([]int, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]int, len(values))
	for i, v := range values {
		vs[i] = int(v)
	}
	return vs, nil
}
func toInt8s(values []int64, e *ae.Error) ([]int8, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]int8, len(values))
	for i, v := range values {
		vs[i] = int8(v)
	}
	return vs, nil
}
func toInt16s(values []int64, e *ae.Error) ([]int16, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]int16, len(values))
	for i, v := range values {
		vs[i] = int16(v)
	}
	return vs, nil
}
func toInt32s(values []int64, e *ae.Error) ([]int32, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]int32, len(values))
	for i, v := range values {
		vs[i] = int32(v)
	}
	return vs, nil
}

func toUints(values []uint64, e *ae.Error) ([]uint, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]uint, len(values))
	for i, v := range values {
		vs[i] = uint(v)
	}
	return vs, nil
}
func toUint8s(values []uint64, e *ae.Error) ([]uint8, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]uint8, len(values))
	for i, v := range values {
		vs[i] = uint8(v)
	}
	return vs, nil
}
func toUint16s(values []uint64, e *ae.Error) ([]uint16, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]uint16, len(values))
	for i, v := range values {
		vs[i] = uint16(v)
	}
	return vs, nil
}
func toUint32s(values []uint64, e *ae.Error) ([]uint32, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]uint32, len(values))
	for i, v := range values {
		vs[i] = uint32(v)
	}
	return vs, nil
}
