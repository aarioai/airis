package httpsvr

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUseRouterUseDoneOrder(t *testing.T) {
	var order []string
	rt := NewRouter()
	rt.UseRouter(func(c *Context) {
		order = append(order, "router")
		c.Next()
	})
	rt.Use(func(c *Context) {
		order = append(order, "use")
		c.Next()
	})
	rt.Done(func(c *Context) {
		order = append(order, "done")
	})
	rt.Get("/ping", func(c *Context) {
		order = append(order, "handler")
		c.WriteString("ok")
		c.Next() // required for Done, same as iris
	})

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Fatalf("response: %d %q", rec.Code, rec.Body.String())
	}
	want := "router,use,handler,done"
	got := strings.Join(order, ",")
	if got != want {
		t.Fatalf("order = %q want %q", got, want)
	}
}

func TestUseRouterRunsOnNotFound(t *testing.T) {
	var ran bool
	rt := NewRouter()
	rt.UseRouter(func(c *Context) {
		ran = true
		c.Header().Set("X-Filter", "1")
		c.Next()
	})
	rt.Get("/ok", func(c *Context) { c.WriteString("ok") })

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/missing", nil))
	if !ran {
		t.Fatal("UseRouter did not run on 404")
	}
	if rec.Header().Get("X-Filter") != "1" {
		t.Fatalf("filter header missing: %v", rec.Header())
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestStopExecutionSkipsDone(t *testing.T) {
	var order []string
	rt := NewRouter()
	rt.Use(func(c *Context) {
		order = append(order, "use")
		c.StopExecution()
		c.StatusCode(http.StatusForbidden)
		c.WriteString("deny")
		c.Next() // should no-op when stopped
	})
	rt.Done(func(c *Context) {
		order = append(order, "done")
	})
	rt.Get("/x", func(c *Context) {
		order = append(order, "handler")
		c.WriteString("ok")
	})

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusForbidden || rec.Body.String() != "deny" {
		t.Fatalf("response: %d %q", rec.Code, rec.Body.String())
	}
	if strings.Join(order, ",") != "use" {
		t.Fatalf("order = %v", order)
	}
}

func TestMultiHandlersAndPartyUse(t *testing.T) {
	var order []string
	rt := NewRouter()
	api := rt.Party("/api")
	api.Use(func(c *Context) {
		order = append(order, "party-use")
		c.Next()
	})
	api.Get("/hi",
		func(c *Context) {
			order = append(order, "h1")
			c.Next()
		},
		func(c *Context) {
			order = append(order, "h2")
			c.WriteString("hi")
		},
	)

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/hi", nil))
	if rec.Body.String() != "hi" {
		t.Fatalf("body %q", rec.Body.String())
	}
	if strings.Join(order, ",") != "party-use,h1,h2" {
		t.Fatalf("order = %v", order)
	}
}

func TestUseOnceReplacesByIdentity(t *testing.T) {
	var n int
	mw := func(c *Context) {
		n++
		c.Next()
	}
	rt := NewRouter()
	rt.UseOnce(mw)
	rt.UseOnce(mw)
	rt.Get("/", func(c *Context) { c.WriteString("ok") })

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if n != 1 {
		t.Fatalf("UseOnce ran %d times", n)
	}
}

func TestAnyAndOnErrorCode(t *testing.T) {
	rt := NewRouter()
	rt.Any("/any", func(c *Context) {
		c.WriteString(c.Method())
	})
	rt.OnErrorCode(http.StatusNotFound, func(c *Context) {
		c.WriteString("missing")
	})

	for _, m := range []string{http.MethodGet, http.MethodPost, http.MethodPut} {
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, httptest.NewRequest(m, "/any", nil))
		if rec.Body.String() != m {
			t.Fatalf("%s: body %q", m, rec.Body.String())
		}
	}

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nope", nil))
	if rec.Code != http.StatusNotFound || rec.Body.String() != "missing" {
		t.Fatalf("404: %d %q", rec.Code, rec.Body.String())
	}
}

func TestPartyFuncAndReset(t *testing.T) {
	var order []string
	rt := NewRouter()
	rt.Use(func(c *Context) {
		order = append(order, "root")
		c.Next()
	})
	rt.PartyFunc("/v1", func(p *Party) {
		p.Reset()
		p.Get("/x", func(c *Context) {
			order = append(order, "handler")
			c.WriteString("x")
		})
	})

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/x", nil))
	if rec.Body.String() != "x" {
		t.Fatalf("body %q", rec.Body.String())
	}
	if strings.Join(order, ",") != "handler" {
		t.Fatalf("Reset should drop parent Use, got %v", order)
	}
}
