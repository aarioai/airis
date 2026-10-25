package httpsvr

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var requestIDSeq atomic.Uint64

// Recover turns a panic into a 500 response and keeps the process alive.
func Recover(c *Context) {
	defer func() {
		if rec := recover(); rec != nil {
			c.StatusCode(http.StatusInternalServerError)
			_, _ = c.WriteString(http.StatusText(http.StatusInternalServerError))
			log.Printf("panic: %v", rec)
		}
	}()
	c.Next()
}

// RequestID stores X-Request-Id, generating one when the client did not send it.
func RequestID(c *Context) {
	id := c.GetHeader("X-Request-Id")
	if id == "" {
		id = newRequestID()
	}
	c.Values().Set(RequestIDKey, id)
	c.Header().Set("X-Request-Id", id)
	c.Next()
}

// AccessLog writes one line after the rest of the chain returns.
func AccessLog(c *Context) {
	start := time.Now()
	c.Next()
	path := ""
	if c.Request() != nil && c.Request().URL != nil {
		path = c.Request().URL.Path
	}
	log.Printf("%s %s %d %s", c.Method(), path, c.Status(), time.Since(start))
}

func newRequestID() string {
	var b [16]byte
	binary.BigEndian.PutUint64(b[:8], uint64(time.Now().UnixNano()))
	binary.BigEndian.PutUint64(b[8:], requestIDSeq.Add(1))
	var out [32]byte
	hex.Encode(out[:], b[:])
	return string(out[:])
}
