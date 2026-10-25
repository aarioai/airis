package acontext

import (
	"net/http"
)

// RequestCarrier is the request surface Bind and ClientIP need.
// httpsvr.Context implements it without acontext importing the router.
type RequestCarrier interface {
	GetHeader(string) string
	RemoteAddr() string
	Request() *http.Request
	ResetRequest(*http.Request)
	ValueString(string) string
	SetValue(string, any)
}
