package request

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/aa/atype/aenum"
	"github.com/aarioai/airis/pkg/afmt"
	"html/template"
	"time"
)

func (r *Request) BodyAudio(p string, required ...bool) (atype.AudioPath, *ae.Error) {
	x, e := r.BodyPath(p, required...)
	if e != nil {
		return "", e
	}
	return atype.AudioPath(x), e
}
func (r *Request) BodyAudios(p string, required ...bool) ([]atype.AudioPath, *ae.Error) {
	paths, e := r.BodyPaths(p, required...)
	if e != nil {
		return nil, e
	}
	audios := make([]atype.AudioPath, len(paths))
	for i, path := range paths {
		audios[i] = atype.AudioPath(path)
	}
	return audios, nil
}
func (r *Request) BodyBooln(p string) (atype.Booln, *ae.Error) {
	b, e := r.BodyBool(p)
	if e != nil {
		return 0, e
	}
	return atype.ToBooln(b), nil
}
func (r *Request) BodyCountry(p string, defaultCountry ...aenum.Country) (aenum.Country, *ae.Error) {
	required := len(defaultCountry) == 0
	cn, e := parseCountry(r.BodyUint16, p, required)
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

func (r *Request) BodyCoordinate(p string, required ...bool) (*atype.Coordinate, *ae.Error) {
	x, e := r.BodyFloat64Map(p, required...)
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
	coord := atype.Coordinate{
		Latitude:  lat,
		Longitude: lng,
		Height:    height,
	}
	return &coord, nil
}
func (r *Request) BodyDate(p string, loc *time.Location, required ...bool) (atype.Date, *ae.Error) {
	x, e := r.Body(p, `^`+atype.DateRegExp+`$`, isRequired(required))
	if e != nil {
		return atype.MinDate, ae.NewBadParam(p)
	}
	return atype.NewDate(x.ReleaseString(), loc), nil
}

func (r *Request) BodyDatetime(p string, loc *time.Location, required ...bool) (atype.Datetime, *ae.Error) {
	x, e := r.Body(p, `^`+atype.DatetimeRegExp+`$`, isRequired(required))
	if e != nil {
		return atype.MinDatetime, ae.NewBadParam(p)
	}
	return atype.NewDatetime(x.ReleaseString(), loc), nil
}

func (r *Request) BodyDecimal(p string, bitSize int, ranges ...atype.Decimal) (atype.Decimal, *ae.Error) {
	return parseDecimal(r.Body, p, bitSize, ranges...)
}
func (r *Request) BodyDist(p string, required ...bool) (atype.Dist, *ae.Error) {
	distri, e := r.BodyDistri(p, required...)
	if e != nil {
		return 0, e
	}
	return distri.Dist(), nil
}
func (r *Request) BodyDistri(p string, required ...bool) (atype.Distri, *ae.Error) {
	x, e := r.BodyUint24(p, isRequired(required))
	return atype.NewDistri(x), e
}
func (r *Request) BodyHtml(p string, required ...any) (template.HTML, *ae.Error) {
	x, e := r.Body(p, required...)
	if e != nil {
		return "", ae.NewBadParam(p)
	}
	return template.HTML(x.ReleaseString()), nil
}
func (r *Request) BodyImage(p string, required ...bool) (atype.ImagePath, *ae.Error) {
	x, e := r.BodyPath(p, required...)
	if e != nil {
		return "", e
	}
	return atype.ImagePath(x), nil
}
func (r *Request) BodyImages(p string, required ...bool) ([]atype.ImagePath, *ae.Error) {
	paths, e := r.BodyPaths(p, required...)
	if e != nil {
		return nil, e
	}
	images := make([]atype.ImagePath, len(paths))
	for i, path := range paths {
		images[i] = atype.ImagePath(path)
	}
	return images, nil
}
func (r *Request) BodyInt24(p string, required ...bool) (atype.Int24, *ae.Error) {
	v, e := parseInt64(r.Body, p, isRequired(required), 24)
	return atype.Int24(v), e
}
func (r *Request) BodyInt24s(p string, required, allowZero bool) ([]atype.Int24, *ae.Error) {
	values, e := r.parseInt64s(r.Body, p, required, allowZero, 24)
	return toInt24s(values, e)
}

func (r *Request) BodyLocation(p string, required ...bool) (*atype.Location, *ae.Error) {
	x, e := r.BodyAnyMap(p, required...)
	if e != nil || x == nil {
		return nil, e
	}
	var loc atype.Location
	lat, ok := x["lat"]
	if !ok {
		e = ae.NewBadParam(p)
		return nil, e
	}

	if loc.Latitude, ok = lat.(float64); !ok {
		e = ae.NewBadParam(p)
		return nil, e
	}
	lng, ok := x["lng"]
	if !ok {
		e = ae.NewBadParam(p)
		return nil, e
	}

	if loc.Longitude, ok = lng.(float64); !ok {
		e = ae.NewBadParam(p)
		return nil, e
	}
	if ht, ok := x["height"]; ok {
		loc.Height, _ = ht.(float64)
	}
	loc.Valid = true
	loc.Name = atype.String(x["name"])
	loc.Address = atype.String(x["address"])
	return &loc, nil
}

func (r *Request) BodyMoney(p string, ranges ...atype.Money) (atype.Money, *ae.Error) {
	return parseMoney(r.Body, p, ranges...)
}

func (r *Request) BodyProvince(p string, required ...bool) (atype.Province, *ae.Error) {
	distri, e := r.BodyDistri(p, required...)
	if e != nil {
		return 0, e
	}
	return distri.Province(), nil
}

func (r *Request) BodySex(p string, xargs ...bool) (aenum.Sex, *ae.Error) {
	return parseSex(r.BodyUint8, p, xargs...)
}
func (r *Request) BodyStatus(p string, xargs ...bool) (aenum.Status, *ae.Error) {
	return parseStatus(r.BodyInt8, p, xargs...)
}

func (r *Request) BodyText(p string, required ...any) (atype.Text, *ae.Error) {
	x, e := r.Body(p, required...)
	if e != nil {
		return "", e
	}
	return atype.NewText(x.ReleaseString(), false), e
}

func (r *Request) BodyUint24(p string, required ...bool) (atype.Uint24, *ae.Error) {
	v, e := parseInt64(r.Body, p, isRequired(required), 24)
	return atype.Uint24(v), e
}

func (r *Request) BodyUint24s(p string, required, allowZero bool) ([]atype.Uint24, *ae.Error) {
	values, e := r.parseUint64s(r.Body, p, required, allowZero, 24)
	return toUint24s(values, e)
}
func (r *Request) BodyVideo(p string, required ...bool) (atype.VideoPath, *ae.Error) {
	x, e := r.BodyPath(p, required...)
	if e != nil {
		return "", e
	}
	return atype.VideoPath(x), e
}
func (r *Request) BodyVideos(p string, required ...bool) ([]atype.VideoPath, *ae.Error) {
	paths, e := r.BodyPaths(p, required...)
	if e != nil {
		return nil, e
	}
	videos := make([]atype.VideoPath, len(paths))
	for i, path := range paths {
		videos[i] = atype.VideoPath(path)
	}
	return videos, nil
}
