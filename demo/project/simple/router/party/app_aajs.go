package party

import (
	"project/simple/app/app_aajs/module/bs/controller"
	"project/simple/router/middleware"

	"github.com/aarioai/airis/aa"
	"github.com/kataras/iris/v12"
)

func RegisterAaJS(app *aa.App, parent *iris.Application, w *middleware.Middleware) {
	p := parent.Party("/")

	bs := controller.New(app)
	p.Get("/ping", bs.Ping)
	p.Get("/ping/redis", bs.PingRedis)
	p.Get("/ping/mysql", bs.PingMySQL)

	p.Get("/health", bs.Health)            // alive health
	p.Get("/health/ready", bs.HealthReady) // service is ready or not

	registerAaJSV1(app, parent, w)
	registerAaJSV1Authed(app, parent, w)
}
