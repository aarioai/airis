package controller

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
    "{{APP_BASE}}/cache"
	"{{APP_BASE}}/module/{{MODULE_NAME}}"
	"{{APP_BASE}}/module/{{MODULE_NAME}}/model"
	"{{APP_BASE}}/service"
	"sync"
	"time"
)

type Controller struct {
	app   *aa.App
	loc   *time.Location
	h     *cache.Cache
	mongo *mongodb.Model
	m     *model.Model
	s     *service.Service
	{{MODULE_NAME}} *{{MODULE_NAME}}.Service
}

var (
	ctrlOnce sync.Once
	ctrlObj  *Controller
)

func New(app *aa.App) *Controller {
	ctrlOnce.Do(func() {
	    s := service.New(app)
		ctrlObj = &Controller{
			app:   app,
			loc:   app.Config.TimeLocation,
			h:     s.Cache(),
			mongo: s.Mongo(),
			m:     model.New(app),
			s:     s,
		    {{MODULE_NAME}}: {{MODULE_NAME}}.New(app),
		}
	})
	return ctrlObj
}

func (c *Controller) mo(t index.Entity) *mongodb.ORMS {
	return c.mongo.ORM(t)
}