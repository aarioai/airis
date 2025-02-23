package controller

import (
	"github.com/aarioai/airis/core"
    "{{APP_BASE}}/cache"
    "{{APP_BASE}}/conf"
	"{{APP_BASE}}/module/{{MODULE_NAME}}"
	"{{APP_BASE}}/module/{{MODULE_NAME}}/model"
	"{{APP_BASE}}/service"
	"{{DRIVER_BASE}}/mongodbhelper"
	"sync"
	"time"
)

type Controller struct {
	app   *core.App
	loc   *time.Location
	h     *cache.Cache
	mongo *mongodbhelper.Model
	m     *model.Model
	s     *service.Service
	{{MODULE_NAME}} *{{MODULE_NAME}}.Service
}

var (
	ctrlOnce sync.Once
	ctrlObj  *Controller
)

func New(app *core.App) *Controller {
	ctrlOnce.Do(func() {
		ctrlObj = &Controller{
			app:   app,
			loc:   app.Config.TimeLocation,
			h:     cache.New(app),
			mongo: mongodbhelper.NewDB(app, conf.MongoDBConfigSection),
			m:     model.New(app),
			s:     service.New(app),
		    {{MODULE_NAME}}: {{MODULE_NAME}}.New(app),
		}
	})
	return ctrlObj
}
