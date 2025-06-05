package atype

import (
	"errors"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
	"strings"
	"time"
)

// IsZeroDate detects a string is empty/zero/min date or datetime
func IsZeroDate(s string) bool {
	if s == "" {
		return true
	}

	if EnableZeroDate {
		if s == ZeroDate.String() || s == ZeroDatetime.String() {
			return true
		}
	}
	return s == MinDate.String() || s == MinDatetime.String()
}

func ToYearMonth(year int, month time.Month) YearMonth {
	if year < 0 {
		return 0
	}
	ym := year*100 + int(month)
	return YearMonth(ym)
}

func ParseYear(year string) (Year, error) {
	y, err := types.ParseUint16(year)
	if err != nil {
		return 0, err
	}
	return Year(y), nil
}
func ParseYearMonth(year, month string) (YearMonth, error) {
	y, err1 := types.ParseUint16(year)
	m, err2 := types.ParseUint8(month)
	if err := ae.FirstError(err1, err2); err != nil {
		return 0, err
	}
	return YearMonth(uint32(y)*100 + uint32(m)), nil
}

func (y Year) String() string      { return types.FormatUint(y) }
func (y YearMonth) String() string { return types.FormatUint(y) }
func (y YMD) String() string       { return types.FormatUint(y) }

// years/months 可为负数
func (y YearMonth) Add(years int, months int, loc *time.Location) YearMonth {
	tm := y.Time(loc).AddDate(years, months, 0)
	return ToYearMonth(tm.Year(), tm.Month())
}
func (y YearMonth) Uin32() uint32 { return uint32(y) }
func (y YearMonth) Date() (int, time.Month) {
	year := int(y) / 100
	month := time.Month(y % 100)
	return year, month
}
func (y YearMonth) Time(loc *time.Location) time.Time {
	year, month := y.Date()
	return time.Date(year, month, 0, 00, 00, 00, 0, loc)
}

func ParseYMD(year, month, day string) (YMD, error) {
	y, err1 := types.ParseUint16(year)
	m, err2 := types.ParseUint8(month)
	d, err3 := types.ParseUint8(day)
	if err := ae.FirstError(err1, err2, err3); err != nil {
		return 0, err
	}
	return YMD(uint(y)*10000 + uint(m)*100 + uint(d)), nil
}
func (y YMD) Uint() uint {
	return uint(y)
}
func (y YMD) Date() Date {
	s := types.FormatUint(y)
	if len(s) != 8 {
		return MinDate
	}
	x := s[0:4] + "-" + s[4:6] + "-" + s[6:]
	return Date(x)
}

// time.Now().In()  loc 直接通过 in 传递
func NewDate(s string, loc *time.Location) Date {
	if IsZeroDate(s) {
		return MinDate
	}
	_, err := time.ParseInLocation("2006-01-02", s, loc)
	if err != nil {
		return MinDate
	}
	return Date(s)
}
func ToDate(t *time.Time) Date {
	if t == nil || t.IsZero() {
		return MinDate
	}
	return Date(t.Format("2006-01-02"))
}
func (d Date) Valid() bool {
	return len(d) == 10
}
func (d Date) String() string { return string(d) }
func (d Date) Time(loc *time.Location) (time.Time, error) {
	if !d.Valid() {
		return time.Time{}, errors.New("min date")
	}
	return time.ParseInLocation("2006-01-02", string(d), loc)
}
func (d Date) OrMin() Date {
	if !d.Valid() {
		return MinDate
	}
	return d
}
func (d Date) OrMax() Date {
	if !d.Valid() {
		return MaxDate
	}
	return d
}
func (d Date) OrNow(loc *time.Location) Date {
	if !d.Valid() {
		return Date(time.Now().In(loc).Format("2006-01-02"))
	}
	return d
}

func (d Date) Ymd() YMD {
	s := strings.ReplaceAll(d.String(), "-", "")
	n, _ := types.ParseUint(s)
	return YMD(n)
}
func (d Date) Int64(loc *time.Location) int64 {
	if !d.Valid() {
		return 0
	}
	tm, err := time.ParseInLocation("2006-01-02", string(d), loc)
	if err != nil {
		return 0
	}
	return tm.Unix()
}
func (d Date) Unix(loc *time.Location) Timestamp {
	if !d.Valid() || IsZeroDate(d.String()) {
		return 0
	}
	t, _ := d.Time(loc)
	return Timestamp(t.Unix())
}

func NewDatetime(s string, loc *time.Location) Datetime {
	if IsZeroDate(s) {
		return MinDatetime
	}
	_, err := time.ParseInLocation("2006-01-02 15:04:05", s, loc)
	if err != nil {
		return MinDatetime
	}
	return Datetime(s)
}
func Now(loc *time.Location) Datetime {
	return ToDatetime(time.Now(), loc)
}

func ToDatetime(t time.Time, loc *time.Location) Datetime {
	if t.IsZero() {
		return MinDatetime
	}
	return Datetime(t.In(loc).Format("2006-01-02 15:04:05"))
}
func ToDatetime2(t *time.Time, loc *time.Location) Datetime {
	if t == nil || t.IsZero() {
		return MinDatetime
	}
	return ToDatetime(*t, loc)
}
func UnixToDatetime(u int64, loc *time.Location) Datetime { return NewTimestamp(u).Datetime(loc) }

func (d Datetime) Valid() bool {
	return len(d) == 19 && d != MinDatetime && d != MaxDatetime && d.String() != "1970-01-01 00:00:00"
}
func (d Datetime) String() string { return string(d) }
func (d Datetime) Time(loc *time.Location) (time.Time, error) {
	if !d.Valid() {
		return time.Time{}, errors.New("min datetime")
	}
	return time.ParseInLocation("2006-01-02 15:04:05", string(d), loc)
}
func (d Datetime) Date() Date { return Date(d[0:len(MinDate)]) }
func (d Datetime) OrMin() Datetime {
	if !d.Valid() {
		return MinDatetime
	}
	return d
}
func (d Datetime) OrMax() Datetime {
	if !d.Valid() {
		return MaxDatetime
	}
	return d
}
func (d Datetime) OrNow(loc *time.Location) Datetime {
	if !d.Valid() {
		return Now(loc)
	}
	return d
}
func (d Datetime) Int64(loc *time.Location) int64 {
	if !d.Valid() {
		return 0
	}
	tm, err := time.ParseInLocation("2006-01-02 15:04:05", string(d), loc)
	if err != nil {
		return 0
	}
	return tm.Unix()
}
func (d Datetime) Unix(loc *time.Location) Timestamp {
	if !d.Valid() {
		return 0
	}
	t, _ := d.Time(loc)
	return Timestamp(t.Unix())
}

func NewTimestamp(u int64) Timestamp {
	return Timestamp(u)
}
func (u Timestamp) Int64() int64 { return int64(u) }
func (u Timestamp) Date(loc *time.Location) Date {
	if u == 0 {
		return MinDate
	}
	t := time.Unix(u.Int64(), 0).In(loc)
	return ToDate(&t)
}
func (u Timestamp) Datetime(loc *time.Location) Datetime {
	if u == 0 {
		return MinDatetime
	}
	t := time.Unix(u.Int64(), 0)
	return ToDatetime(t, loc)
}

func (u Timestamp) Add(duration Second) Timestamp {
	return u + Timestamp(duration)
}

func (u Timestamp) Sub(b Timestamp) Second {
	return Second(u - b)
}

func (u Timestamp) SubUnix(b int64) Second {
	return Second(int64(u) - b)
}
