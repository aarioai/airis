package atype

import (
	"github.com/aarioai/airis/pkg/basic"
	"strings"
)

// Iris path parameter types
// See https://iris-go.gitbook.io/iris/contents/routing/routing-path-parameter-types

const (
	InvalidWeekday Weekday = -1
	Sunday         Weekday = 0
	Monday         Weekday = 1
	Tuesday        Weekday = 2
	Wednesday      Weekday = 3
	Thursday       Weekday = 4
	Friday         Weekday = 5
	Saturday       Weekday = 6
)

var WeekdayName = map[Weekday]string{
	InvalidWeekday: "",
	Sunday:         "sunday",
	Monday:         "monday",
	Tuesday:        "tuesday",
	Wednesday:      "wednesday",
	Thursday:       "thursday",
	Friday:         "friday",
	Saturday:       "saturday",
}

// ToWeekday Converts weekday code and English names to weekday code
// English names (ignore cases):
//
//	sunday, monday, tuesday, wednesday, thursday, friday, saturday
//	sun., mon., tue., wed., thu., fri., sat.,
//	sun, mon, tue, wed, thu, fri, sat,
//	su., mo., tu., we., th., fr., sa.,
//	su, mo, tu, we, th, fr, sa,
//	tues., tues, thur., thur, thurs., thurs
//
// Return [0-6] from sunday to saturday, -1 to invalid weekday
//
// Example
//
//	ToWeekday(0)  // Returns 0, 0 is sunday
//	ToWeekday(-3) // Returns -1, -1 is an invalid week day
//	ToWeekday('Sunday') // Returns 0
//	ToWeekday('FRI.')  // Returns 5
func ToWeekday(s string) Weekday {
	if s == "" {
		return InvalidWeekday
	}
	if len(s) == 1 {
		return basic.Ter(s[0] >= '0' && s[0] <= '6', Weekday(s[0]-'0'), InvalidWeekday)
	}

	s = strings.ReplaceAll(strings.Trim(s, " "), ".", "")
	switch strings.ToLower(s) {
	case "0", "sunday", "sun", "su":
		return Sunday
	case "1", "monday", "mon", "mo":
		return Monday
	case "2", "tuesday", "tue", "tu", "tues":
		return Tuesday
	case "3", "wednesday", "wed", "we":
		return Wednesday
	case "4", "thursday", "thur", "th", "thurs":
		return Thursday
	case "5", "friday", "fri", "fr":
		return Friday
	case "6", "saturday", "sat", "sa":
		return Saturday
	default:
		return InvalidWeekday
	}
}

func NewWeekday(n uint8) Weekday {
	w := Weekday(n)
	if w.Valid() {
		return w
	}
	return InvalidWeekday
}
func (s Weekday) Valid() bool {
	return s >= Sunday && s <= Saturday
}
