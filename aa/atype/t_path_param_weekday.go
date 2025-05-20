package atype

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
