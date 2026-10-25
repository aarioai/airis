package httpsvr

import (
		"strings"
)

// HeaderWild reads a header from the raw name, the canonical name, and an X- prefix.
func HeaderWild(ctx *Context, key string) string {
	value := ctx.GetHeader(key)
	if value != "" {
		return value
	}
	newKey := titleHeaderKey(key)
	if newKey != key {
		if value = ctx.GetHeader(newKey); value != "" {
			return value
		}
	}
	if !strings.HasPrefix(newKey, "X-") {
		return ctx.GetHeader("X-" + newKey)
	}
	return ""
}

// QueryWild reads a value from the path, the query string, a header, and optionally a cookie.
func QueryWild(ctx *Context, key string, includeCookie bool) string {
	value := ctx.Params().GetString(key)
	if value != "" {
		return value
	}
	if ctx.Request() != nil {
		if value = ctx.Request().URL.Query().Get(key); value != "" {
			return value
		}
	}
	if value = HeaderWild(ctx, key); value != "" {
		return value
	}
	if includeCookie {
		return ctx.GetCookie(key)
	}
	return ""
}
