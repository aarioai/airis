package controller

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
    "project/microservice/app/infra/cache"
	"project/microservice/app/infra/module/task"
	"project/microservice/app/infra/module/task/model"
	"project/microservice/app/infra/service"
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
	task *task.Service
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
		    task: task.New(app),
		}
	})
	return ctrlObj
}

func (c *Controller) mo(t index.Entity) *mongodb.ORMS {
	return c.mongo.ORM(t)
}