package httpsvr

import (
	"net/http"

	"github.com/aarioai/airis/aa/ae"
)

func JSON(ctx *Context, data any) error {
	d := Body{
		Code: ae.OK,
		Msg:  "OK",
		Data: data,
	}
	return ctx.JSON(d)
}

func JsonOK(ctx *Context) error {
	return JSON(ctx, nil)
}

func JsonCode(ctx *Context, code int) error {
	d := Body{
		Code: code,
		Msg:  ae.Text(code),
		Data: nil,
	}
	return ctx.JSON(d)
}

func JsonError(ctx *Context, code int, msg string) error {
	d := Body{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	return ctx.JSON(d)
}

func JsonE(ctx *Context, e *ae.Error) error {
	if e == nil {
		return JsonOK(ctx)
	}
	d := Body{
		Code: e.Code,
		Msg:  e.Msg,
		Data: nil,
	}
	return ctx.JSON(d)
}

func StatusHandler(status int) Handler {
	return func(ctx *Context) {
		defer ctx.Next()
		ctx.StatusCode(status)
	}
}

func WriteHandler[T string | []byte](msg T) Handler {
	return func(ctx *Context) {
		defer ctx.Next()
		_, _ = ctx.Write([]byte(msg))
	}
}

func StatusMessageHandler[T string | []byte](status int, msg T) Handler {
	return func(ctx *Context) {
		defer ctx.Next()
		if status > 0 && status != http.StatusOK {
			ctx.StatusCode(status)
		}
		_, _ = ctx.Write([]byte(msg))
	}
}
