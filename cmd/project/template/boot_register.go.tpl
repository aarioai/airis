package boot

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/kataras/iris/v12"
)

// Register global handlers
func register(app *aa.App) {
	response.RegisterGlobalServeContentTypes([]string{"application/json", "text/html"})

	response.RegisterGlobalErrorHandler(func(ictx iris.Context, r *request.Request, contentType string, d response.Body) (int, error, bool) {
		if d.Code == ae.NotModified {
			ictx.StatusCode(ae.NotModified)
			return 0, nil, false
		}

		dbg := ictx.Values().GetBoolDefault("_debug", false)
		if d.Code >= 500 && !dbg {
			e := ae.New(d.Code, d.Msg).WithCaller(3)
			d.Msg = "Internal Server Error" // hide server errmsg
			alog.Error(e.Text())
		}

		if contentType == request.CtHTML.String() {
			ictx.ViewData("code", d.Code)
			ictx.ViewData("msg", d.Msg)
			ictx.View("error/error.html")
		}

		return 0, nil, true
	})
}
