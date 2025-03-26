package response

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/kataras/iris/v12"
)

func JSON(ictx iris.Context, data any, opts ...iris.JSON) error {
	d := Body{
		Code: ae.OK,
		Msg:  "OK",
		Data: data,
	}
	return ictx.JSON(d, opts...)
}

func JsonOK(ictx iris.Context, opts ...iris.JSON) error {
	return JSON(ictx, nil, opts...)
}

func JsonCode(ictx iris.Context, code int, opts ...iris.JSON) error {
	d := Body{
		Code: code,
		Msg:  ae.Text(code),
		Data: nil,
	}
	return ictx.JSON(d, opts...)
}
func JsonError(ictx iris.Context, code int, msg string, opts ...iris.JSON) error {
	d := Body{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	return ictx.JSON(d, opts...)
}
func JsonE(ictx iris.Context, e *ae.Error, opts ...iris.JSON) error {
	if e == nil {
		return JsonOK(ictx)
	}
	d := Body{
		Code: e.Code,
		Msg:  e.Msg,
		Data: nil,
	}
	return ictx.JSON(d, opts...)
}

func StatusHandler(status int) func(iris.Context) {
	return func(ictx iris.Context) {
		defer ictx.Next()
		ictx.StatusCode(status)
	}
}
func WriteHandler[T string | []byte](msg T) func(iris.Context) {
	return func(ictx iris.Context) {
		defer ictx.Next()
		_, _ = ictx.Write([]byte(msg))
	}
}

func Handler[T string | []byte](status int, msg T) func(iris.Context) {
	return func(ictx iris.Context) {
		defer ictx.Next()
		if status > 0 && status != iris.StatusOK {
			ictx.StatusCode(status)
		}
		_, _ = ictx.Write([]byte(msg))
	}
}
