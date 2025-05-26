package atype

const (
	False            Booln    = 0
	True             Booln    = 1
	ZeroDate         Date     = "0000-00-00"
	ZeroDatetime     Datetime = "0000-00-00 00:00:00"
	MysqlMinDate     Date     = "1000-01-01"
	MysqlMaxDate     Date     = "9999-12-31"
	MysqlMinDatetime Datetime = "1000-01-01 00:00:00"
	MysqlMaxDatetime Datetime = "9999-12-31 23:59:59"

	DateRegExp     = `([12]\d{3}-[01]\d-[0-3]\d)|(0000-00-00)|(9999-12-31)`
	DatetimeRegExp = `([12]\d{3}-[01]\d-[0-3]\d\s[0-2]\d:[0-5]\d:[0-5]\d)|(0000-00-00\s00:00:00)|(9999-12-31\s23:59:59)`
)

var (
	EnableZeroDate = true
	MinDate        = MysqlMinDate
	MaxDate        = MysqlMaxDate
	MinDatetime    = MysqlMinDatetime
	MaxDatetime    = MysqlMaxDatetime
)
