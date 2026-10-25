# httpsvr

Native HTTP context and router for airis. Built on `net/http` only.

Request parsing (`Query*` / `Body*`) and DTO responses (`TryWrite` / `WriteE` / …) live on `*httpsvr.Context`.

## Benchmarks

```bash
go test ./aa/httpsvr -bench=. -benchmem -count=1
```

## Stress tests

```bash
# Live TCP server, 64 workers for 3s
go test ./aa/httpsvr -run TestStressHTTPServer -count=1 -v

# In-process ServeHTTP fan-out
go test ./aa/httpsvr -run TestStressServeHTTPParallel -count=1 -v
```

`-short` skips the stress tests.
