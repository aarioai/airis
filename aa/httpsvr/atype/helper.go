package atype

import (
	"github.com/aarioai/airis/aa/ae"
	datatype "github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/aa/atype/aenum"
	"github.com/aarioai/airis/pkg/types"
)

func isRequired(requireds []bool) bool {
	return len(requireds) == 0 || requireds[0]
}

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
	cn, ok := aenum.NewCountry(sts)
	if !ok {
		return 0, ae.NewBadParam(p)
	}
	return cn, nil
}

func parseSex(method func(string, ...bool) (uint8, *ae.Error), p string, xargs ...bool) (aenum.Sex, *ae.Error) {
	required := len(xargs) == 0 || xargs[0]
	sts, e := method(p, required)
	if e != nil {
		return 0, e
	}
	return aenum.NewSex(sts), nil
}

func parseDecimalString(s, p string, bitSize int, ranges ...datatype.Decimal) (datatype.Decimal, *ae.Error) {
	if s == "" {
		return 0, nil
	}
	x, err := types.ParseInt64(s, bitSize)
	if err != nil {
		return 0, ae.NewBadParam(p)
	}
	n := datatype.Decimal(x)
	rangeMin := datatype.MinDecimal
	rangeMax := datatype.MaxDecimal
	if len(ranges) > 0 {
		rangeMin = ranges[0]
	}
	if len(ranges) == 2 {
		rangeMax = ranges[1]
	}
	if n < rangeMin || n > rangeMax {
		return 0, ae.NewBadParam(p)
	}
	return n, nil
}

func parseMoneyString(s, p string, ranges ...datatype.Money) (datatype.Money, *ae.Error) {
	if s == "" {
		return 0, nil
	}
	x, err := types.ParseInt64(s, 64)
	if err != nil {
		return 0, ae.NewBadParam(p)
	}
	n := datatype.Money(x)
	rangeMin := datatype.MinMoney
	rangeMax := datatype.MaxMoney
	if len(ranges) > 0 {
		rangeMin = ranges[0]
	}
	if len(ranges) == 2 {
		rangeMax = ranges[1]
	}
	if n < rangeMin || n > rangeMax {
		return 0, ae.NewBadParam(p)
	}
	return n, nil
}

func toInt24s(values []int64, e *ae.Error) ([]datatype.Int24, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}
	vs := make([]datatype.Int24, len(values))
	for i, v := range values {
		vs[i] = datatype.Int24(v)
	}
	return vs, nil
}

func toUint24s(values []uint64, e *ae.Error) ([]datatype.Uint24, *ae.Error) {
	if e != nil || len(values) == 0 {
		return nil, e
	}
	vs := make([]datatype.Uint24, len(values))
	for i, v := range values {
		vs[i] = datatype.Uint24(v)
	}
	return vs, nil
}
