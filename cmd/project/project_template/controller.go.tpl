package controller

import (
	"github.com/aarioai/airis/core"
	"{{APP_BASE}}/module/{{MODULE_NAME}}"
	"{{APP_BASE}}/service"
	"sync"
	"time"
)

type Controller struct {
	app *core.App
	loc *time.Location
	s   *service.Service
	{{MODULE_NAME}} *{{MODULE_NAME}}.Service
}

var (
	ctrlOnce sync.Once
	ctrlObj  *Controller
)

func New(app *core.App) *Controller {
	ctrlOnce.Do(func() {
		ctrlObj = &Controller{
		    app: app,
		    loc: app.Config.TimeLocation,
		    s: service.New(app),
		    {{MODULE_NAME}}: {{MODULE_NAME}}.New(app),
		}
	})
	return ctrlObj
}
