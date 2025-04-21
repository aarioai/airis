package atype

import "time"

type DurationSeconds int64

func ToDurationSeconds(t time.Duration) DurationSeconds {
	return DurationSeconds(t.Seconds())
}

func (d DurationSeconds) Duration() time.Duration {
	return time.Duration(d) * time.Second
}
