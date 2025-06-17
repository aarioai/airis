package ss

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
	"project/simple/app/app_aajs/cache"
	"project/simple/app/app_aajs/module/ss/model"
	"project/simple/app/app_aajs/private"
	"project/simple/app/app_aajs/service"
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
