package atype

import (
	"html/template"
	"time"

	"github.com/aarioai/airis/aa/ae"
	datatype "github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/aa/atype/aenum"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/types"
)

// Body groups Body parsers that return atype / aenum values.
// Use as ctx.BodyAtype.AudioPaths("files").
type Body struct{ Src BodySource }

func (a Body) AudioPath(p string, required ...bool) (datatype.AudioPath, *ae.Error) {
	x, e := a.Src.BodyPath(p, required...)
	if e != nil {
		return "", e
	}
	return datatype.AudioPath(x), e
}

func (a Body) AudioPaths(p string, required ...bool) ([]datatype.AudioPath, *ae.Error) {
	paths, e := a.Src.BodyPaths(p, required...)
	if e != nil {
		return nil, e
	}
	audios := make([]datatype.AudioPath, len(paths))
	for i, path := range paths {
		audios[i] = datatype.AudioPath(path)
	}
	return audios, nil
}

func (a Body) Booln(p string) (datatype.Booln, *ae.Error) {
	b, e := a.Src.BodyBool(p)
	if e != nil {
		return 0, e
	}
	return datatype.ToBooln(b), nil
}

func (a Body) Country(p string, defaultCountry ...aenum.Country) (aenum.Country, *ae.Error) {
	required := len(defaultCountry) == 0
	cn, e := parseCountry(a.Src.BodyUint16, p, required)
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

func (a Body) Coordinate(p string, required ...bool) (*datatype.Coordinate, *ae.Error) {
	x, e := a.Src.BodyFloat64Map(p, required...)
	if e != nil || x == nil {
		return nil, e
	}
	lat, ok := x["lat"]
	if !ok {
		return nil, ae.NewBadParam(p)
	}
	lng, ok := x["lng"]
	if !ok {
		return nil, ae.NewBadParam(p)
	}
	height, _ := x["height"]
	coord := datatype.Coordinate{
		Latitude:  lat,
		Longitude: lng,
		Height:    height,
	}
	return &coord, nil
}

func (a Body) Date(p string, loc *time.Location, required ...bool) (datatype.Date, *ae.Error) {
	x, e := a.Src.BodyString(p, `^`+datatype.DateRegExp+`$`, isRequired(required))
	if e != nil {
		return datatype.MinDate, ae.NewBadParam(p)
	}
	return datatype.NewDate(x, loc), nil
}

func (a Body) Datetime(p string, loc *time.Location, required ...bool) (datatype.Datetime, *ae.Error) {
	x, e := a.Src.BodyString(p, `^`+datatype.DatetimeRegExp+`$`, isRequired(required))
	if e != nil {
		return datatype.MinDatetime, ae.NewBadParam(p)
	}
	return datatype.NewDatetime(x, loc), nil
}

func (a Body) Decimal(p string, bitSize int, ranges ...datatype.Decimal) (datatype.Decimal, *ae.Error) {
	s, e := a.Src.BodyString(p, false)
	if e != nil {
		return 0, e
	}
	return parseDecimalString(s, p, bitSize, ranges...)
}

func (a Body) Dist(p string, required ...bool) (datatype.Dist, *ae.Error) {
	distri, e := a.Distri(p, required...)
	if e != nil {
		return 0, e
	}
	return distri.Dist(), nil
}

func (a Body) Distri(p string, required ...bool) (datatype.Distri, *ae.Error) {
	x, e := a.Uint24(p, isRequired(required))
	return datatype.NewDistri(x), e
}

func (a Body) Html(p string, required ...any) (template.HTML, *ae.Error) {
	x, e := a.Src.BodyString(p, required...)
	if e != nil {
		return "", ae.NewBadParam(p)
	}
	return template.HTML(x), nil
}

func (a Body) ImagePath(p string, required ...bool) (datatype.ImagePath, *ae.Error) {
	x, e := a.Src.BodyPath(p, required...)
	if e != nil {
		return "", e
	}
	return datatype.ImagePath(x), nil
}

func (a Body) ImagePaths(p string, required ...bool) ([]datatype.ImagePath, *ae.Error) {
	paths, e := a.Src.BodyPaths(p, required...)
	if e != nil {
		return nil, e
	}
	images := make([]datatype.ImagePath, len(paths))
	for i, path := range paths {
		images[i] = datatype.ImagePath(path)
	}
	return images, nil
}

func (a Body) Int24(p string, required ...bool) (datatype.Int24, *ae.Error) {
	s, e := a.Src.BodyString(p, isRequired(required))
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

func (a Body) Int24s(p string, required, allowZero bool) ([]datatype.Int24, *ae.Error) {
	values, e := a.Src.BodyInt64s(p, required, allowZero)
	return toInt24s(values, e)
}

func (a Body) Location(p string, required ...bool) (*datatype.Location, *ae.Error) {
	x, e := a.Src.BodyAnyMap(p, required...)
	if e != nil || x == nil {
		return nil, e
	}
	var loc datatype.Location
	lat, ok := x["lat"]
	if !ok {
		return nil, ae.NewBadParam(p)
	}
	if loc.Latitude, ok = lat.(float64); !ok {
		return nil, ae.NewBadParam(p)
	}
	lng, ok := x["lng"]
	if !ok {
		return nil, ae.NewBadParam(p)
	}
	if loc.Longitude, ok = lng.(float64); !ok {
		return nil, ae.NewBadParam(p)
	}
	if ht, ok := x["height"]; ok {
		loc.Height, _ = ht.(float64)
	}
	loc.Valid = true
	loc.Name = datatype.String(x["name"])
	loc.Address = datatype.String(x["address"])
	return &loc, nil
}

func (a Body) Money(p string, ranges ...datatype.Money) (datatype.Money, *ae.Error) {
	s, e := a.Src.BodyString(p, false)
	if e != nil {
		return 0, e
	}
	return parseMoneyString(s, p, ranges...)
}

func (a Body) Province(p string, required ...bool) (datatype.Province, *ae.Error) {
	distri, e := a.Distri(p, required...)
	if e != nil {
		return 0, e
	}
	return distri.Province(), nil
}

func (a Body) Sex(p string, xargs ...bool) (aenum.Sex, *ae.Error) {
	return parseSex(a.Src.BodyUint8, p, xargs...)
}

func (a Body) Status(p string, xargs ...bool) (aenum.Status, *ae.Error) {
	return parseStatus(a.Src.BodyInt8, p, xargs...)
}

func (a Body) Text(p string, required ...any) (datatype.Text, *ae.Error) {
	x, e := a.Src.BodyString(p, required...)
	if e != nil {
		return "", e
	}
	return datatype.NewText(x, false), e
}

func (a Body) Uint24(p string, required ...bool) (datatype.Uint24, *ae.Error) {
	s, e := a.Src.BodyString(p, isRequired(required))
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

func (a Body) Uint24s(p string, required, allowZero bool) ([]datatype.Uint24, *ae.Error) {
	values, e := a.Src.BodyUint64s(p, required, allowZero)
	return toUint24s(values, e)
}

func (a Body) VideoPath(p string, required ...bool) (datatype.VideoPath, *ae.Error) {
	x, e := a.Src.BodyPath(p, required...)
	if e != nil {
		return "", e
	}
	return datatype.VideoPath(x), e
}

func (a Body) VideoPaths(p string, required ...bool) ([]datatype.VideoPath, *ae.Error) {
	paths, e := a.Src.BodyPaths(p, required...)
	if e != nil {
		return nil, e
	}
	videos := make([]datatype.VideoPath, len(paths))
	for i, path := range paths {
		videos[i] = datatype.VideoPath(path)
	}
	return videos, nil
}
