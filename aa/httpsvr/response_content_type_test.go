package httpsvr

import "testing"

func TestNegotiateContentType(t *testing.T) {
	serve := []string{"application/json", "text/html"}
	cases := []struct {
		header string
		want   string
	}{
		{"application/json", "application/json"},
		{"text/html, application/json;q=0.8", "text/html"},
		{"text/plain;q=0.1, application/json;q=0.9", "application/json"},
		{"text/*;q=0.2, application/json", "application/json"},
		{"*/*", "application/json"},
		{"application/json; charset=utf-8", "application/json"},
		{"text/plain", ""},
	}
	for _, c := range cases {
		if got := negotiateContentType(c.header, serve); got != c.want {
			t.Errorf("negotiate(%q) = %q, want %q", c.header, got, c.want)
		}
	}
}
