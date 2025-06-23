package middleware

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/kataras/iris/v12"
)

func (w *Middleware) Auth(ictx iris.Context) {
	auth := ictx.GetHeader("Authorization")
	if auth != "Bearer helloworld" {
		response.JsonE(ictx, ae.ErrorUnauthorized)
		return
	}
	ictx.Next()
}
