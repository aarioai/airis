package router

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/helpers/debug"
)

func Serve(app *aa.App, prof *debug.Profile) {
	serveHTTP(app, "app_aajs")
}
