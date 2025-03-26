package response

import (
	"github.com/kataras/iris/v12"
	"net/http"
)

func JsonOK(ictx iris.Context, data any, opts ...iris.JSON) error {
	d := Body{
		Code: http.StatusOK,
		Msg:  "OK",
		Data: data,
	}
	return ictx.JSON(d, opts...)
}

func JsonE(ictx iris.Context, code int, msg string, opts ...iris.JSON) error {
	d := Body{
		Code: code,
		Msg:  msg,
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
