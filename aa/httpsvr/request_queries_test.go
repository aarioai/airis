package httpsvr

import (
	"net/http"
	"testing"
)

func TestQueriesMergeInjectedOverUser(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.test/users?page=1&q=go", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-From", "client")
	c := &Context{
		r:               req,
		injectedQueries: map[string]any{"page": "9", "uid": "3"},
		injectedHeaders: map[string]any{"X-From": "injected", "X-Debug": "1"},
		paramsInjected:  true,
	}

	got := c.Queries()
	if got["page"] != "9" || got["q"] != "go" || got["uid"] != "3" {
		t.Fatalf("Queries() = %#v", got)
	}
	if _, ok := got["X-From"]; ok {
		t.Fatalf("Queries() included header %#v", got)
	}

	headers := c.Headers()
	if headers["X-From"] != "injected" || headers["X-Debug"] != "1" {
		t.Fatalf("Headers() = %#v", headers)
	}
	if _, ok := headers["page"]; ok {
		t.Fatalf("Headers() included query %#v", headers)
	}
}

func TestQueriesWithoutInjection(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.test/users?page=2", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &Context{r: req, paramsInjected: true}
	got := c.Queries()
	if got["page"] != "2" || len(got) != 1 {
		t.Fatalf("Queries() = %#v", got)
	}
}
