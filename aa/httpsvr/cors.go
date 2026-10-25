package httpsvr

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSOptions configures the CORS middleware.
type CORSOptions struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	MaxAge           int
	AllowCredentials bool
}

// CORS sets cross-origin headers and answers preflight requests.
// A wildcard origin together with credentials echoes the request Origin instead of "*".
func CORS(opt CORSOptions) Handler {
	methods := strings.Join(opt.AllowedMethods, ", ")
	headers := strings.Join(opt.AllowedHeaders, ", ")
	exposed := strings.Join(opt.ExposedHeaders, ", ")
	return func(c *Context) {
		origin := c.GetHeader("Origin")
		allow := allowOrigin(opt.AllowedOrigins, origin, opt.AllowCredentials)
		if allow != "" {
			c.Header().Set("Access-Control-Allow-Origin", allow)
			if allow != "*" {
				c.Header().Add("Vary", "Origin")
			}
		}
		if opt.AllowCredentials && allow != "" && allow != "*" {
			c.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if exposed != "" {
			c.Header().Set("Access-Control-Expose-Headers", exposed)
		}
		if c.Method() == http.MethodOptions {
			if methods != "" {
				c.Header().Set("Access-Control-Allow-Methods", methods)
			}
			if headers != "" {
				reqHeaders := c.GetHeader("Access-Control-Request-Headers")
				if headers == "*" && reqHeaders != "" {
					c.Header().Set("Access-Control-Allow-Headers", reqHeaders)
				} else {
					c.Header().Set("Access-Control-Allow-Headers", headers)
				}
			}
			if opt.MaxAge > 0 {
				c.Header().Set("Access-Control-Max-Age", strconv.Itoa(opt.MaxAge))
			}
			c.StatusCode(http.StatusNoContent)
			c.FlushStatus()
			return
		}
		c.Next()
	}
}

func allowOrigin(allowed []string, origin string, credentials bool) string {
	for _, item := range allowed {
		if item == "*" {
			if credentials && origin != "" {
				return origin
			}
			return "*"
		}
		if item == origin {
			return origin
		}
	}
	return ""
}
