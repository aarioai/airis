package response

import (
	"github.com/kataras/iris/v12"
)

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
