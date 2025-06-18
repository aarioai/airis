package task

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
	"project/microservice/app/app_luexu/cache"
	"project/microservice/app/app_luexu/module/task/model"
	"project/microservice/app/app_luexu/private"
	"project/microservice/app/app_luexu/service"
	"sync"
	"time"
)

type Service struct {
	app   *aa.App
	loc   *time.Location
	private *private.Service
	h     *cache.Cache
	mongo *mongodb.Model
	m     *model.Model
	s     *service.Service
}

var (
	svcOnce sync.Once
	svcObj  *Service
)

func New(app *aa.App) *Service {
	svcOnce.Do(func() {
		s := service.New(app)
		svcObj = &Service{
			app:   app,
			loc:   app.Config.TimeLocation,
			private: private.New(app),
			h:     s.Cache(),
			mongo: s.Mongo(),
			m:     model.New(app),
			s:     s,
		}
	})
	return svcObj
}

func (s *Service) mo(t index.Entity) *mongodb.ORMS {
	return s.mongo.ORM(t)
}
