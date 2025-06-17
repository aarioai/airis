package party

import (
	"github.com/aarioai/airis/aa"
	"github.com/kataras/iris/v12"
	"project/simple/app/app_aajs/module/bs/controller"
	"project/simple/router/middleware"
)

func registerAaJSV1(app *aa.App, parent *iris.Application, w *middleware.Middleware) {
	p := parent.Party("/test/v1")

	bs := controller.New(app)
	p.Get("/ping", bs.GetPing)
}
