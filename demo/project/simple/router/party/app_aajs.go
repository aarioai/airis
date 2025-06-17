package party

import (
	"github.com/aarioai/airis/aa"
	"github.com/kataras/iris/v12"
	"project/simple/app/app_aajs/module/bs/controller"
	"project/simple/router/middleware"
)

func RegisterAaJS(app *aa.App, parent *iris.Application, w *middleware.Middleware) {
	p := parent.Party("/")

	bs := controller.New(app)
	p.Get("/ping", bs.GetPing)

	registerAaJSV1(app, parent, w)
}
