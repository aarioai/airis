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

	p.Get("/users", bs.GetUsersWithPaging)
	p.Get("/users/page/{page:uint8}", bs.GetUsersWithPaging)
	p.Get("/users/sex/{sex:uint8}", bs.GetUsersWithPaging)
	p.Get("/users/sex/{sex:uint8}/page/{page:uint8}", bs.GetUsersWithPaging)
	p.Get("/users/{uid:uint64}", bs.GetUser)

	p.Head("/restful", bs.HeadRestful)
	p.Post("/restful", bs.PostRestful)
	p.Put("/restful/{id:int}", bs.PutRestful)
	p.Patch("/restful", bs.PatchRestful)
	p.Get("/restful/{response:string}", bs.GetRestful)
	p.Delete("/restful", bs.DeleteRestful)

	p.Post("/login", bs.PostLogin)                 // password login
	p.Head("/auth/access_token", bs.HeadUserToken) // detect access token validate
	p.Put("/auth/access_token", bs.GrantUserToken) // refresh access token
}
