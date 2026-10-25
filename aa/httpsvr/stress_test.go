package httpsvr

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestStressHTTPServer hits a live http.Server with concurrent clients.
// Run: go test ./aa/httpsvr -run TestStressHTTPServer -count=1 -v
// Skip under -short.
func TestStressHTTPServer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	rt := NewRouter()
	rt.UseRouter(Recover, RequestID)
	p := rt.Party("/v1")
	p.Get("/ping", func(c *Context) {
		_, _ = c.WriteString("PONG")
	})
	p.Get("/users/{uid:uint64}", func(c *Context) {
		_ = c.JSON(map[string]any{
			"uid": c.Params().Store[0].ValueRaw,
		})
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{Handler: rt}
	go func() { _ = srv.Serve(ln) }()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	base := "http://" + ln.Addr().String()
	const (
		workers  = 64
		duration = 3 * time.Second
	)

	var (
		okCount   atomic.Int64
		errCount  atomic.Int64
		byteCount atomic.Int64
	)

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        workers * 2,
			MaxIdleConnsPerHost: workers * 2,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	deadline := time.Now().Add(duration)
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func(id int) {
			defer wg.Done()
			i := 0
			for time.Now().Before(deadline) {
				url := base + "/v1/ping"
				if i%3 == 0 {
					url = fmt.Sprintf("%s/v1/users/%d", base, id*1000+i)
				}
				i++
				resp, err := client.Get(url)
				if err != nil {
					errCount.Add(1)
					continue
				}
				n, _ := io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					okCount.Add(1)
					byteCount.Add(n)
				} else {
					errCount.Add(1)
				}
			}
		}(w)
	}
	wg.Wait()

	ok := okCount.Load()
	errs := errCount.Load()
	elapsed := duration.Seconds()
	rps := float64(ok) / elapsed
	t.Logf("workers=%d duration=%s ok=%d err=%d bytes=%d rps=%.0f",
		workers, duration, ok, errs, byteCount.Load(), rps)
	if ok == 0 {
		t.Fatal("no successful requests")
	}
	if errs > ok/100 {
		t.Fatalf("error rate too high: ok=%d err=%d", ok, errs)
	}
}

// TestStressServeHTTPParallel stress-tests ServeHTTP without a real TCP stack.
func TestStressServeHTTPParallel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	rt := NewRouter()
	rt.UseRouter(Recover, RequestID)
	rt.Party("/v1").Get("/users/{uid:uint64}", func(c *Context) {
		_, _ = c.WriteString("ok")
	})

	const (
		workers   = 128
		perWorker = 5000
	)
	var failed atomic.Int64
	var wg sync.WaitGroup
	wg.Add(workers)
	start := time.Now()
	for w := 0; w < workers; w++ {
		go func(id int) {
			defer wg.Done()
			req := httptestRequest(fmt.Sprintf("/v1/users/%d", id+1))
			for i := 0; i < perWorker; i++ {
				rec := newRecorder()
				rt.ServeHTTP(rec, req)
				if rec.code != http.StatusOK {
					failed.Add(1)
				}
			}
		}(w)
	}
	wg.Wait()
	elapsed := time.Since(start)
	total := int64(workers * perWorker)
	rps := float64(total) / elapsed.Seconds()
	t.Logf("in-process: workers=%d total=%d failed=%d elapsed=%s rps=%.0f",
		workers, total, failed.Load(), elapsed, rps)
	if failed.Load() != 0 {
		t.Fatalf("failed=%d", failed.Load())
	}
}

type stubResponse struct {
	header http.Header
	code   int
	body   []byte
}

func newRecorder() *stubResponse {
	return &stubResponse{header: make(http.Header), code: http.StatusOK}
}

func (s *stubResponse) Header() http.Header { return s.header }
func (s *stubResponse) Write(b []byte) (int, error) {
	if s.code == 0 {
		s.code = http.StatusOK
	}
	s.body = append(s.body, b...)
	return len(b), nil
}
func (s *stubResponse) WriteHeader(code int) { s.code = code }

func httptestRequest(path string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	return req
}
