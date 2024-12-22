package atype

import (
	"github.com/aarioai/airis/pkg/types"
	"strings"
)

type SepPercents string
type SepMoneys string

func ToSepPercents(elems []Decimal) SepPercents {
	// strings.Concat 类同
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return SepPercents(types.FormatInt(elems[0]))
	}
	deli := ","
	n := (len(elems) - 1) + (len(elems) * types.MaxInt16Len)
	var b strings.Builder
	b.Grow(n)
	b.WriteString(types.FormatInt(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(deli)
		b.WriteString(types.FormatInt(s))
	}

	return SepPercents(b.String())
}

func (t SepPercents) Percents() []Decimal {
	if t == "" {
		return nil
	}
	arr := strings.Split(string(t), ",")
	v := make([]Decimal, len(arr))
	for i, a := range arr {
		p, err := types.ParseInt(a)
		if err == nil {
			v[i] = Decimal(p)
		}
	}
	return v
}

func ToSepMoneys(elems []Money) SepMoneys {
	// strings.Concat 类同
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return SepMoneys(types.FormatInt(elems[0]))
	}
	deli := ","
	n := (len(elems) - 1) + (len(elems) * types.MaxInt64Len)
	var b strings.Builder
	b.Grow(n)
	b.WriteString(types.FormatInt(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(deli)
		b.WriteString(types.FormatInt(s))
	}

	return SepMoneys(b.String())
}
func (t SepMoneys) Moneys() []Money {
	if t == "" {
		return nil
	}
	arr := strings.Split(string(t), ",")
	v := make([]Money, len(arr))
	for i, a := range arr {
		p, err := types.ParseInt(a)
		if err == nil {
			v[i] = Money(p)
		}
	}
	return v
}
