package httpsvr

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func benchStack() *Router {
	settings := NewSettings()
	rt := NewRouter()
	rt.UseHooks(settings)
	rt.UseRouter(Recover, RequestID)
	p := rt.Party("/v1")
	p.Get("/ping", func(c *Context) {
		_, _ = c.WriteString("PONG")
	})
	p.Get("/echo/{uid:uint64}", func(c *Context) {
		defer c.Next()
		uid, e := c.QueryUint64("uid")
		if e != nil {
			_, _ = c.WriteE(e)
			return
		}
		_, _ = c.WriteData(map[string]any{"uid": uid})
	})
	p.Post("/login", func(c *Context) {
		defer c.Next()
		account, e0 := c.BodyString("account", false)
		password, e1 := c.BodyString("password", false)
		if e := c.FirstError(e0, e1); e != nil {
			return
		}
		_, _ = c.WriteData(map[string]string{"account": account, "password": password})
	})
	return rt
}

func BenchmarkHTTPXPing(b *testing.B) {
	rt := benchStack()
	req := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)
	}
}

func BenchmarkHTTPXJSONEcho(b *testing.B) {
	rt := benchStack()
	req := httptest.NewRequest(http.MethodGet, "/v1/echo/42", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)
	}
}

func BenchmarkHTTPXJSONEchoParallel(b *testing.B) {
	rt := benchStack()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest(http.MethodGet, "/v1/echo/99", nil)
		for pb.Next() {
			rec := httptest.NewRecorder()
			rt.ServeHTTP(rec, req)
		}
	})
}

func BenchmarkHTTPXFormLogin(b *testing.B) {
	rt := benchStack()
	body := "account=10001&password=secret"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)
	}
}
