package httpsvr

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTypedPathAndMiddleware(t *testing.T) {
	rt := NewRouter()
	rt.UseHooks("app-a")
	rt.UseRouter(Recover, RequestID)
	p := rt.Party("/v1")
	p.Get("/users/{uid:uint64}", func(c *Context) {
		if c.Hooks() != "app-a" {
			t.Errorf("hooks = %v", c.Hooks())
		}
		if c.Params().GetString("uid") != "42" {
			t.Errorf("uid = %v", c.Params().Store)
		}
		if _, ok := c.Params().Store[0].ValueRaw.(uint64); !ok {
			t.Errorf("uid type %T", c.Params().Store[0].ValueRaw)
		}
		if c.Values().GetString(RequestIDKey) == "" {
			t.Error("missing request id")
		}
		c.WriteString("ok")
	})
	p.Party("/authed", func(c *Context) {
		if c.GetHeader("Authorization") == "" {
			c.StatusCode(http.StatusUnauthorized)
			c.WriteString("no")
			return
		}
		c.Next()
	}).Get("/ping", func(c *Context) {
		c.WriteString("pong")
	})

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/users/42", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Fatalf("users: %d %q", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/users/nope", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("bad uid status %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/authed/ping", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("auth status %d", rec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/authed/ping", nil)
	req.Header.Set("Authorization", "Bearer helloworld")
	rec = httptest.NewRecorder()
	rt.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || rec.Body.String() != "pong" {
		t.Fatalf("authed: %d %q", rec.Code, rec.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	rt := NewRouter()
	rt.Party("/v1").Get("/ping", func(c *Context) { c.WriteString("ok") })
	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/ping", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Allow"), http.MethodGet) {
		t.Fatalf("Allow=%q", rec.Header().Get("Allow"))
	}
}

func TestHooksAreNotShared(t *testing.T) {
	a := NewRouter()
	b := NewRouter()
	a.UseHooks("a")
	b.UseHooks("b")
	a.Party("/").Get("/", func(c *Context) {
		c.WriteString(c.Hooks().(string))
	})
	b.Party("/").Get("/", func(c *Context) {
		c.WriteString(c.Hooks().(string))
	})
	ra := httptest.NewRecorder()
	a.ServeHTTP(ra, httptest.NewRequest(http.MethodGet, "/", nil))
	rb := httptest.NewRecorder()
	b.ServeHTTP(rb, httptest.NewRequest(http.MethodGet, "/", nil))
	if ra.Body.String() != "a" || rb.Body.String() != "b" {
		t.Fatalf("a=%q b=%q", ra.Body.String(), rb.Body.String())
	}
}
