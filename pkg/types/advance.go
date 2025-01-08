package types

import (
	"log"
	"strings"
	"time"
)

// ParseDuration parses a duration string.
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// A duration string is a possibly signed sequence of decimal numbers,
// each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
func ParseDuration(s string, defaultValue ...time.Duration) time.Duration {
	var dv time.Duration
	if len(defaultValue) > 0 {
		dv = defaultValue[0]
	}
	if s == "" {
		return dv
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Println("[error] invalid time duration: " + s)
		return dv
	}
	return d
}

// StringifyDuration stringifies time duration to a duration string
func StringifyDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
