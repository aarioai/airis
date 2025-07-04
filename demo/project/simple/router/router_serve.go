package router

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"project/simple/router/middleware"
	"project/simple/router/party"
)

func serveHTTP(app *aa.App, serviceName string) {
	p := iris.Default()
	p.Logger().SetLevel("info")

	p.UseRouter(recover.New(), logger.New())
	go func() {
		<-app.GlobalContext.Done()
		alog.Stop(serviceName)
		p.Shutdown(app.GlobalContext)
	}()

	w := middleware.New(app)
	party.RegisterAaJS(app, p, w)

	port := app.Config.GetString(serviceName + ".port")
	serve := iris.Addr(":" + port)
	iwc := iris.WithConfiguration(iris.Configuration{
		LogLevel:          "info",
		DisableStartupLog: true,
		// DisableInterruptHandler:           false,
		// DisablePathCorrection:             false,
		// EnablePathEscape:                  false,
		// FireMethodNotAllowed:              false,
		// DisableBodyConsumptionOnUnmarshal: false,
		// DisableAutoFireStatusCode:         false,
		TimeFormat:         app.Config.TimeFormat,
		Charset:            "UTF-8",
		IgnoreServerErrors: []string{iris.ErrServerClosed.Error()},
	})
	ae.PanicOnErrs(p.Run(serve, iwc))
}
