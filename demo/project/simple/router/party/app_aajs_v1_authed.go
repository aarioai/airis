package party

import (
	"github.com/aarioai/airis/aa"
	"github.com/kataras/iris/v12"
	"project/simple/app/app_aajs/module/bs/controller"
	"project/simple/router/middleware"
)

func registerAaJSV1Authed(app *aa.App, parent *iris.Application, w *middleware.Middleware) {
	p := parent.Party("/v1/authed")

	bs := controller.New(app)

	p.Get("/users", bs.GetUsersWithPaging)
	p.Get("/users/page/{page:uint8}", bs.GetUsersWithPaging)
	p.Get("/users/{uid:uint64}", bs.GetUser)
	p.Get("/users/sex/{sex:uint8}", bs.GetUsersWithPaging)
	p.Get("/users/sex/{sex:uint8}/page/{page:uint8}", bs.GetUsersWithPaging)

	p.Get("/users2", bs.GetUsers)
	p.Get("/users2/uid/{uid:string}", bs.GetUsers)

	p.Post("/users", bs.PostUser)
	p.Put("/users/{uid:uint64}", bs.PutUser)
	p.Patch("/users/{uid:uint64}", bs.PatchUser)
	p.Delete("/users/{uid:uint64}", bs.DeleteUser)
}
