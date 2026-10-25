package atype

import (
	"time"

	"github.com/aarioai/airis/aa/ae"
	datatype "github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/aa/atype/aenum"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/types"
)

const (
	paramPage     = "page"
	paramPageSize = "page_size"
	paramPageEnd  = "page_end"
)

// Query groups Query parsers that return atype / aenum values.
// Use as ctx.QueryAtype.Sex("sex").
type Query struct{ Src QuerySource }

func (a Query) Paging() datatype.Paging {
	page, _ := a.Src.QueryUint(paramPage, false)
	pageSize, _ := a.Src.QueryUint8(paramPageSize, false)
	pageEnd, _ := a.Src.QueryUint(paramPageEnd, false)
	return datatype.NewPaging(page, pageEnd, pageSize)
}

func (a Query) PagingWithSize(maxPageSize uint8) datatype.Paging {
	page, _ := a.Src.QueryUint(paramPage, false)
	pageSize, _ := a.Src.QueryUint8(paramPageSize, false)
	pageEnd, _ := a.Src.QueryUint(paramPageEnd, false)
	if maxPageSize > 0 && (pageSize == 0 || pageSize > maxPageSize) {
		pageSize = maxPageSize
	}
	return datatype.NewPaging(page, pageEnd, pageSize)
}

func (a Query) PagingWithFirstPageEnd(defaultPageEnd uint) datatype.Paging {
	page, _ := a.Src.QueryUint(paramPage, false)
	pageSize, _ := a.Src.QueryUint8(paramPageSize, false)
	pageEnd, _ := a.Src.QueryUint(paramPageEnd, false)
	if pageSize > datatype.DefaultPageSize {
		pageSize = datatype.DefaultPageSize
	}
	if page < 2 && pageEnd <= page {
		pageEnd = defaultPageEnd
	}
	return datatype.NewPaging(page, pageEnd, pageSize)
}

func (a Query) Booln(p string) (datatype.Booln, *ae.Error) {
	b, e := a.Src.QueryBool(p)
	if e != nil {
		return 0, e
	}
	return datatype.ToBooln(b), nil
}

func (a Query) Country(p string, defaultCountry ...aenum.Country) (aenum.Country, *ae.Error) {
	required := len(defaultCountry) == 0
	cn, e := parseCountry(a.Src.QueryUint16, p, required)
	if e != nil {
		return cn, e
	}
	if len(defaultCountry) == 0 {
		return cn, nil
	}
	if cn == 0 {
		return afmt.First(defaultCountry), nil
	}
	return cn, nil
}

func (a Query) Date(p string, loc *time.Location, required ...bool) (datatype.Date, *ae.Error) {
	x, e := a.Src.QueryString(p, `^`+datatype.DateRegExp+`$`, isRequired(required))
	if e != nil {
		return datatype.MinDate, ae.NewBadParam(p)
	}
	return datatype.NewDate(x, loc), nil
}

func (a Query) Datetime(p string, loc *time.Location, required ...bool) (datatype.Datetime, *ae.Error) {
	x, e := a.Src.QueryString(p, `^`+datatype.DatetimeRegExp+`$`, isRequired(required))
	if e != nil {
		return datatype.MinDatetime, ae.NewBadParam(p)
	}
	return datatype.NewDatetime(x, loc), nil
}

func (a Query) Decimal(p string, bitSize int, ranges ...datatype.Decimal) (datatype.Decimal, *ae.Error) {
	s, e := a.Src.QueryString(p, false)
	if e != nil {
		return 0, e
	}
	return parseDecimalString(s, p, bitSize, ranges...)
}

func (a Query) Dist(p string, required ...bool) (datatype.Dist, *ae.Error) {
	distri, e := a.Distri(p, required...)
	if e != nil {
		return 0, e
	}
	return distri.Dist(), nil
}

func (a Query) Distri(p string, required ...bool) (datatype.Distri, *ae.Error) {
	x, e := a.Uint24(p, isRequired(required))
	return datatype.NewDistri(x), e
}

func (a Query) Int24(p string, required ...bool) (datatype.Int24, *ae.Error) {
	s, e := a.Src.QueryString(p, isRequired(required))
	if e != nil {
		return 0, e
	}
	if s == "" {
		return 0, nil
	}
	v, err := types.ParseInt64(s, 24)
	if err != nil {
		return 0, ae.NewBadParam(p)
	}
	return datatype.Int24(v), nil
}

func (a Query) Int24s(p string, required, allowZero bool) ([]datatype.Int24, *ae.Error) {
	values, e := a.Src.QueryInt64s(p, required, allowZero)
	return toInt24s(values, e)
}

func (a Query) Money(p string, ranges ...datatype.Money) (datatype.Money, *ae.Error) {
	s, e := a.Src.QueryString(p, false)
	if e != nil {
		return 0, e
	}
	return parseMoneyString(s, p, ranges...)
}

func (a Query) Province(p string, required ...bool) (datatype.Province, *ae.Error) {
	distri, e := a.Distri(p, required...)
	if e != nil {
		return 0, e
	}
	return distri.Province(), nil
}

func (a Query) Sex(p string, xargs ...bool) (aenum.Sex, *ae.Error) {
	return parseSex(a.Src.QueryUint8, p, xargs...)
}

func (a Query) Status(p string, xargs ...bool) (aenum.Status, *ae.Error) {
	return parseStatus(a.Src.QueryInt8, p, xargs...)
}

func (a Query) StatusEnum(p string, enums ...aenum.Status) (aenum.Status, *ae.Error) {
	status, e := parseStatus(a.Src.QueryInt8, p)
	if e != nil || len(enums) == 0 {
		return status, e
	}
	for _, x := range enums {
		if x == status {
			return x, nil
		}
	}
	return status, ae.NewBadParam(p, "check status enum")
}

func (a Query) Uint24(p string, required ...bool) (datatype.Uint24, *ae.Error) {
	s, e := a.Src.QueryString(p, isRequired(required))
	if e != nil {
		return 0, e
	}
	if s == "" {
		return 0, nil
	}
	v, err := types.ParseUint64(s, 24)
	if err != nil {
		return 0, ae.NewBadParam(p)
	}
	return datatype.Uint24(v), nil
}

func (a Query) Uint24s(p string, required, allowZero bool) ([]datatype.Uint24, *ae.Error) {
	values, e := a.Src.QueryUint64s(p, required, allowZero)
	return toUint24s(values, e)
}
