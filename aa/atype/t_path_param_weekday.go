package atype

import "strings"

// Iris path parameter types
// See https://iris-go.gitbook.io/iris/contents/routing/routing-path-parameter-types

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

var WeekdayName = map[Weekday]string{
	Sunday:    "sunday",
	Monday:    "monday",
	Tuesday:   "tuesday",
	Wednesday: "wednesday",
	Thursday:  "thursday",
	Friday:    "friday",
	Saturday:  "saturday",
}

func NewWeekday(n uint8) (Weekday, bool) {
	w := Weekday(n)
	return w, w <= Saturday
}
func (s Weekday) Valid() bool {
	return s <= Saturday
}

func normalizeWeekDay(s string) string {
	if s == "" {
		return ""
	}
	if len(s) == 1 {
		r := s[0]
		if r >= '0' && r <= '6' {
			return s
		}
		return ""
	}
	s = strings.ToLower(strings.Replace(s, ".", "", 1))
	switch s {
	case "sunday", "sun":
		return "0"
	case "monday", "mon":
		return "1"
	case "tuesday", "tue", "tues":
		return "2"
	case "wednesday", "wed":
		return "3"
	case "thursday", "thu", "thur", "thurs":
		return "4"
	case "friday", "fri":
		return "5"
	case "saturday", "sat":
		return "6"
	}
	return ""
}

func NewWeekDay(s string) (WeekDay, bool) {
	s = normalizeWeekDay(s)
	w := WeekDay(s)
	return w, s != ""
}

func (w WeekDay) Valid() bool {
	return normalizeWeekDay(string(w)) != ""
}

func (w WeekDay) Weekday() (Weekday, bool) {
	s := normalizeWeekDay(string(w))
	if s == "" {
		return 0, false
	}
	return NewWeekday(s[0] - '0' + '0')
}
