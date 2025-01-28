package request

import (
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/aenum"
	"github.com/aarioai/airis/core/atype"
	"github.com/aarioai/airis/pkg/types"
	"reflect"
	"regexp"
	"strings"
)

func parseStatus(method func(string, ...bool) (int8, *ae.Error), p string, xargs ...bool) (aenum.Status, *ae.Error) {
	required := len(xargs) == 0 || xargs[0]
	sts, e := method(p, required)
	if e != nil {
		return 0, e
	}

	status, ok := aenum.NewStatus(sts)
	if !ok {
		return 0, ae.NewBadParam(p)
	}
	return status, nil
}
func parseCountry(method func(string, ...bool) (uint16, *ae.Error), p string, xargs ...bool) (aenum.Country, *ae.Error) {
	required := len(xargs) == 0 || xargs[0]
	sts, e := method(p, required)
	if e != nil {
		return 0, e
	}

	status, ok := aenum.NewCountry(sts)
	if !ok {
		return 0, ae.NewBadParam(p)
	}
	return status, nil
}
func parseDecimal(method func(string, ...any) (*RawValue, *ae.Error), p string, bitSize int, ranges ...atype.Decimal) (atype.Decimal, *ae.Error) {
	x, e := parseInt64(method, p, false, bitSize)
	n := atype.Decimal(x)
	rangeMin := atype.MinDecimal
	rangeMax := atype.MaxDecimal
	if len(ranges) > 0 {
		rangeMin = ranges[0]
	}
	if len(ranges) == 2 {
		rangeMax = ranges[1]
	}
	if n < rangeMin || n > rangeMax {
		return 0, ae.NewBadParam(p)
	}
	return n, e
}
func parseMoney(method func(string, ...any) (*RawValue, *ae.Error), p string, ranges ...atype.Money) (atype.Money, *ae.Error) {
	x, e := parseInt64(method, p, false, 64)
	n := atype.Money(x)
	rangeMin := atype.MinMoney
	rangeMax := atype.MaxMoney
	if len(ranges) > 0 {
		rangeMin = ranges[0]
	}
	if len(ranges) == 2 {
		rangeMax = ranges[1]
	}
	if n < rangeMin || n > rangeMax {
		return 0, ae.NewBadParam(p)
	}
	return n, e
}
func parseInt64(method func(string, ...any) (*RawValue, *ae.Error), p string, required bool, bitSize int) (int64, *ae.Error) {
	value, e := method(p, required)
	if e != nil {
		return 0, e
	}
	v, err := types.ParseInt64(value.String(), bitSize)
	if err != nil {
		return 0, ae.NewBadParam(p)
	}
	return v, nil
}
func parseSex(method func(string, ...bool) (uint8, *ae.Error), p string, xargs ...bool) (aenum.Sex, *ae.Error) {
	required := len(xargs) == 0 || xargs[0]
	sts, e := method(p, required)
	if e != nil {
		return 0, e
	}

	return aenum.NewSex(sts), nil
}
func parseUint64(method func(string, ...any) (*RawValue, *ae.Error), p string, required bool, bitSize int) (uint64, *ae.Error) {
	value, e := method(p, required)
	if e != nil {
		return 0, e
	}
	v, err := types.ParseUint64(value.String())
	if err != nil {
		return 0, ae.NewBadParam(p)
	}
	return v, nil
}
func (r *Request) parseInt64s(method func(string, ...any) (*RawValue, *ae.Error), p string, required, allowZero bool, bitSize int) ([]int64, *ae.Error) {
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
		return r.sepInt64s(method, p, ",", required, allowZero, bitSize)
	}
	if reflect.TypeOf(d).Kind() != reflect.Slice {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	// 有可能是 [1,"2",3] 这种混合的数组
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
func (r *Request) parseUint64s(method func(string, ...any) (*RawValue, *ae.Error), p string, required, allowZero bool, bitSize int) ([]uint64, *ae.Error) {
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
		return r.separatedUint64s(method, p, ",", required, allowZero, bitSize)
	}
	if reflect.TypeOf(d).Kind() != reflect.Slice {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	// 有可能是 [1,"2",3] 这种混合的数组
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
func (r *Request) parseStrings(method func(string, ...any) (*RawValue, *ae.Error), p string, required, allowEmptyString bool) ([]string, *ae.Error) {
	q, e := r.Body(p, required)
	if e != nil {
		return nil, e
	}
	defer q.Release()
	var v []string
	d := q.Raw()
	if d == nil {
		if required {
			return nil, ae.NewBadParam(p)
		}
		return nil, nil
	}
	if _, ok := d.(string); ok {
		return sepStrings(method, p, ",", required, allowEmptyString)
	}
	if d != nil {
		switch reflect.TypeOf(d).Kind() {
		case reflect.Slice: // 有可能是 [1,"2",3] 这种混合的数组
			s := reflect.ValueOf(d)
			v = make([]string, 0, s.Len())
			for i := 0; i < s.Len(); i++ {
				// 不能用 s.Index(i).String()，否则返回：<interface {} Value>
				ts := atype.String(s.Index(i).Interface())
				if allowEmptyString || (ts != "" && ts != "0") {
					v = append(v, ts)
				}
			}
		}
	}
	if len(v) == 0 && required {
		return nil, ae.NewBadParam(p)
	}
	return v, nil
}

// @warn 禁止传递 float
func sepStrings(method func(string, ...any) (*RawValue, *ae.Error), p string, sep string, required, allowEmptyString bool) ([]string, *ae.Error) {
	s, e := method(p, required)

	if e != nil {
		return nil, e
	}
	// 将换行符当作切割符号
	re := regexp.MustCompile(`\s*[\r\n]+\s*`)
	x := re.ReplaceAllString(s.String(), sep)
	arr := strings.Split(x, sep)
	b := make([]string, 0)
	for _, a := range arr {
		if allowEmptyString {
			b = append(b, a)
		} else {
			a = strings.Trim(a, " ")
			if a != "" {
				b = append(b, a)
			}
		}

	}
	if len(b) == 0 && required {
		return nil, ae.NewBadParam(p)
	}
	return b, nil
}

// @warn 禁止传递 float
func (r *Request) sepInt64s(method func(string, ...any) (*RawValue, *ae.Error), p string, sep string, required, allowZero bool, bitSize int) ([]int64, *ae.Error) {
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
func (r *Request) separatedUint64s(method func(string, ...any) (*RawValue, *ae.Error), p string, sep string, required, allowZero bool, bitSize int) ([]uint64, *ae.Error) {
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
func toInt24s(values []int64, e *ae.Error) ([]atype.Int24, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]atype.Int24, len(values))
	for i, v := range values {
		vs[i] = atype.Int24(v)
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
func toUint24s(values []uint64, e *ae.Error) ([]atype.Uint24, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}

	vs := make([]atype.Uint24, len(values))
	for i, v := range values {
		vs[i] = atype.Uint24(v)
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
