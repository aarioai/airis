package request

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/afmt"
	"time"
)

func (r *Request) Paging() atype.Paging {
	page, _ := r.QueryUint(ParamPage, false)
	pageSize, _ := r.QueryUint8(ParamPageSize, false)
	pageEnd, _ := r.QueryUint(ParamPageEnd, false)
	return atype.NewPaging(page, pageEnd, pageSize)
}

func (r *Request) PagingWithSize(maxPageSize uint8) atype.Paging {
	page, _ := r.QueryUint(ParamPage, false)
	pageSize, _ := r.QueryUint8(ParamPageSize, false)
	pageEnd, _ := r.QueryUint(ParamPageEnd, false)
	if maxPageSize > 0 && (pageSize == 0 || pageSize > maxPageSize) {
		pageSize = maxPageSize
	}
	return atype.NewPaging(page, pageEnd, pageSize)
}

func (r *Request) PagingWithFirstPageEnd(defaultPageEnd uint) atype.Paging {
	page, _ := r.QueryUint(ParamPage, false)
	pageSize, _ := r.QueryUint8(ParamPageSize, false)
	pageEnd, _ := r.QueryUint(ParamPageEnd, false)
	if pageSize > atype.DefaultPageSize {
		pageSize = atype.DefaultPageSize
	}
	if page < 2 && pageEnd <= page {
		pageEnd = defaultPageEnd
	}
	return atype.NewPaging(page, pageEnd, pageSize)
}

func (r *Request) QueryBooln(p string) (atype.Booln, *ae.Error) {
	b, e := r.QueryBool(p)
	if e != nil {
		return 0, e
	}
	return atype.ToBooln(b), nil
}

func (r *Request) QueryCountry(p string, defaultCountry ...aenum.Country) (aenum.Country, *ae.Error) {
	required := len(defaultCountry) == 0
	cn, e := parseCountry(r.QueryUint16, p, required)
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

func (r *Request) QueryDate(p string, loc *time.Location, required ...bool) (atype.Date, *ae.Error) {
	x, e := r.Query(p, `^`+aenum.DateRegExp+`$`, isRequired(required))
	if e != nil {
		return atype.MinDate, ae.NewBadParam(p)
	}
	return atype.NewDate(x.ReleaseString(), loc), nil
}

func (r *Request) QueryDatetime(p string, loc *time.Location, required ...bool) (atype.Datetime, *ae.Error) {
	x, e := r.Query(p, `^`+aenum.DatetimeRegExp+`$`, isRequired(required))
	if e != nil {
		return atype.MinDatetime, ae.NewBadParam(p)
	}
	return atype.NewDatetime(x.ReleaseString(), loc), nil
}

func (r *Request) QueryDecimal(p string, bitSize int, ranges ...atype.Decimal) (atype.Decimal, *ae.Error) {
	return parseDecimal(r.Query, p, bitSize, ranges...)
}

func (r *Request) QueryDist(p string, required ...bool) (atype.Dist, *ae.Error) {
	distri, e := r.QueryDistri(p, required...)
	if e != nil {
		return 0, e
	}
	return distri.Dist(), nil
}
func (r *Request) QueryDistri(p string, required ...bool) (atype.Distri, *ae.Error) {
	x, e := r.QueryUint24(p, isRequired(required))
	return atype.NewDistri(x), e
}

func (r *Request) QueryInt24(p string, required ...bool) (atype.Int24, *ae.Error) {
	v, e := parseInt64(r.Query, p, isRequired(required), 24)
	return atype.Int24(v), e
}
func (r *Request) QueryInt24s(p string, required, allowZero bool) ([]atype.Int24, *ae.Error) {
	values, e := r.parseInt64s(r.Query, p, required, allowZero, 24)
	return toInt24s(values, e)
}
func (r *Request) QueryMoney(p string, ranges ...atype.Money) (atype.Money, *ae.Error) {
	return parseMoney(r.Query, p, ranges...)
}

func (r *Request) QueryProvince(p string, required ...bool) (atype.Province, *ae.Error) {
	distri, e := r.QueryDistri(p, required...)
	if e != nil {
		return 0, e
	}
	return distri.Province(), nil
}

func (r *Request) QuerySex(p string, xargs ...bool) (aenum.Sex, *ae.Error) {
	return parseSex(r.QueryUint8, p, xargs...)
}
func (r *Request) QueryStatus(p string, xargs ...bool) (aenum.Status, *ae.Error) {
	return parseStatus(r.QueryInt8, p, xargs...)
}
func (r *Request) QueryUint24(p string, required ...bool) (atype.Uint24, *ae.Error) {
	v, e := parseInt64(r.Query, p, isRequired(required), 24)
	return atype.Uint24(v), e
}
func (r *Request) QueryUint24s(p string, required, allowZero bool) ([]atype.Uint24, *ae.Error) {
	values, e := r.parseUint64s(r.Query, p, required, allowZero, 24)
	return toUint24s(values, e)
}
