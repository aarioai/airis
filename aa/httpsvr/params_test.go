package httpsvr

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParamsTypedGetters(t *testing.T) {
	rt := NewRouter()
	rt.Get("/u/{uid:uint64}/n/{n:int}/b/{ok:bool}", func(c *Context) {
		uid, ok := c.Params().GetUint64("uid")
		if !ok || uid != 7 {
			t.Fatalf("uid=%v ok=%v", uid, ok)
		}
		n, ok := c.Params().GetInt("n")
		if !ok || n != -3 {
			t.Fatalf("n=%v ok=%v", n, ok)
		}
		b, ok := c.Params().GetBool("ok")
		if !ok || !b {
			t.Fatalf("bool=%v ok=%v", b, ok)
		}
		e, ok := c.Params().GetEntry("uid")
		if !ok || e.Key != "uid" {
			t.Fatalf("entry=%v", e)
		}
		c.WriteString("ok")
	})

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/u/7/n/-3/b/true", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Fatalf("%d %q", rec.Code, rec.Body.String())
	}
}
