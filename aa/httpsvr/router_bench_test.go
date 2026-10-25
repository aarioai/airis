package httpsvr

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func benchRouter() *Router {
	rt := NewRouter()
	rt.UseRouter(Recover, RequestID)
	p := rt.Party("/v1")
	p.Get("/ping", func(c *Context) {
		_, _ = c.WriteString("PONG")
	})
	p.Get("/users/{uid:uint64}", func(c *Context) {
		_, _ = c.WriteString(c.Params().GetString("uid"))
	})
	p.Get("/users/sex/{sex:uint8}/page/{page:uint8}", func(c *Context) {
		_, _ = c.WriteString(c.Params().GetString("page"))
	})
	p.Get("/items/{id:int}", func(c *Context) {
		_ = c.JSON(map[string]any{"id": c.Params().Store[0].ValueRaw})
	})
	return rt
}

func BenchmarkRouterStatic(b *testing.B) {
	rt := benchRouter()
	req := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)
	}
}

func BenchmarkRouterTypedUint64(b *testing.B) {
	rt := benchRouter()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/123456789", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)
	}
}

func BenchmarkRouterMultiParam(b *testing.B) {
	rt := benchRouter()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/sex/1/page/3", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)
	}
}

func BenchmarkRouterJSON(b *testing.B) {
	rt := benchRouter()
	req := httptest.NewRequest(http.MethodGet, "/v1/items/42", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)
	}
}

func BenchmarkRouterParallel(b *testing.B) {
	rt := benchRouter()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest(http.MethodGet, "/v1/users/99", nil)
		for pb.Next() {
			rec := httptest.NewRecorder()
			rt.ServeHTTP(rec, req)
		}
	})
}

func BenchmarkRouterManyRoutes(b *testing.B) {
	rt := NewRouter()
	p := rt.Party("/api")
	for i := 0; i < 200; i++ {
		path := "/r" + strconv.Itoa(i) + "/{id:uint64}"
		p.Get(path, func(c *Context) {
			_, _ = c.WriteString("ok")
		})
	}
	req := httptest.NewRequest(http.MethodGet, "/api/r199/7", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)
	}
}

func BenchmarkStdlibServeMux(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "PONG")
	})
	req := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
	}
}
