package {{PACKAGE_NAME}}

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
    "{{APP_BASE}}/cache"
    "{{APP_BASE}}/conf"
	"sync"
	"time"
)

type Service struct {
	app   *aa.App
	loc   *time.Location
	h     *cache.Cache
	mongo *mongodb.Model
}

var (
	svcOnce sync.Once
	svcObj  *Service
)

func New(app *aa.App) *Service {
	svcOnce.Do(func() {
		svcObj = &Service{
			app:   app,
			loc:   app.Config.TimeLocation,
            h:     cache.New(app),
            mongo: mongodb.NewDB(app, conf.MongoCfgSection),
		}
	})
	return svcObj
}

func (s *Service) mo(t index.Entity) *mongodb.ORMS {
	return s.mongo.ORM(t)
}