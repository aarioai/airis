package consumer

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
	"project/simple/app/app_aajs/cache"
    "project/simple/app/app_aajs/module/bs"
    "project/simple/app/app_aajs/module/cms"
    "project/simple/app/app_aajs/module/ss"
    "project/simple/app/app_aajs/private"
	"project/simple/app/app_aajs/service"
	"sync"
	"time"
)

type Service struct {
	app   *aa.App
	loc   *time.Location
	h     *cache.Cache
	mongo *mongodb.Model
	s     *service.Service
    bs      *bs.Service
    cms     *cms.Service
    ss      *ss.Service
    private *private.Service
}

var (
	svcOnce sync.Once
	svcObj  *Service
)

func New(app *aa.App) *Service {
	svcOnce.Do(func() {
		s := service.New(app)
		svcObj = &Service{
			app:     app,
			loc:     app.Config.TimeLocation,
			h:       s.Cache(),
			mongo:   s.Mongo(),
			s:       s,
			bs:      bs.New(app),
			cms:     cms.New(app),
			ss:      ss.New(app),
			private: private.New(app),
		}
	})
	return svcObj
}

func (s *Service) mo(t index.Entity) *mongodb.ORMS {
	return s.mongo.ORM(t)
}