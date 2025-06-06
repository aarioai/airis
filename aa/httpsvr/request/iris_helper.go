package request

import (
	"github.com/kataras/iris/v12"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

// HeaderWild read http header wildly (including standard or non-standard format)
//  1. origin format, e.g. name, Name, user_agent
//  2. standard format, e.g. Referer, User-Agent
//  3. self-defined format, i.e. starts with X-, e.g. X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
//
// Suggest any non-standard header should be allowed by the web server cors AllowedMethods
func HeaderWild(ictx iris.Context, key string) string {
	// 1. origin format, e.g. key, Name, user_agent
	value := ictx.GetHeader(key)
	if value != "" {
		return value
	}
	// 2. standard format, e.g. Referer, User-Agent
	newKey := cases.Title(language.English).String(strings.ReplaceAll(key, "_", "-"))
	if newKey != key {
		if value = ictx.GetHeader(newKey); value != "" {
			return value
		}
	}
	// 3. self-defined format, i.e. starts with X-, e.g. X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
	if !strings.HasPrefix(newKey, "X-") {
		return ictx.GetHeader("X-" + newKey)
	}
	return ""
}

// QueryWild query parameter from URL parameter, URL query, header and cookie
// Example  csrf_token: in url params? -> Csrf-Token: in header?  X-Csrf-Token: in header-> csrf_token: in cookie
func QueryWild(ictx iris.Context, key string, includeCookie bool) string {
	// 1. query from URL parameter
	value := ictx.Params().GetString(key)
	if value != "" {
		return value
	}
	// 2. query from URL query string
	if value = ictx.Request().URL.Query().Get(key); value != "" {
		return value
	}
	// 3. query from header
	if value = HeaderWild(ictx, key); value != "" {
		return value
	}
	// 4. query from Cookie
	if includeCookie {
		return ictx.GetCookie(key)
	}
	return ""
}
