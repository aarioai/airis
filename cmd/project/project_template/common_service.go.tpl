package {{PACKAGE_NAME}}

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
	"{{APP_BASE}}/cache"
	"{{APP_BASE}}/service"
	"sync"
	"time"
)

type Service struct {
	app   *aa.App
	loc   *time.Location
	h     *cache.Cache
	mongo *mongodb.Model
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
			h:     s.Cache(),
			mongo: s.Mongo(),
			s:     s,
		}
	})
	return svcObj
}

func (s *Service) mo(t index.Entity) *mongodb.ORMS {
	return s.mongo.ORM(t)
}