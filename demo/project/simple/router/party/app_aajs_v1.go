package party

import (
	"github.com/aarioai/airis/aa"
	"github.com/kataras/iris/v12"
	"project/simple/app/app_aajs/module/bs/controller"
	"project/simple/router/middleware"
)

func registerAaJSV1(app *aa.App, parent *iris.Application, w *middleware.Middleware) {
	p := parent.Party("/v1")

	bs := controller.New(app)
	p.Get("/ping", bs.GetPing)

	p.Get("/users", bs.GetUsers)
	p.Get("/users/page/{page:uint8}", bs.GetUsers)
	p.Get("/users/{uid:uint64}", bs.GetUser)
	p.Post("/users", bs.PostUser)
	p.Put("/users/{uid:uint64}", bs.PutUser)
	p.Patch("/users/{uid:uint64}", bs.PatchUser)
	p.Delete("/users/{uid:uint64}", bs.DeleteUser)
}
