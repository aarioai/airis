package {{PACKAGE_NAME}}

import (
    "github.com/aarioai/airis-driver/driver/mongodbhelper"
	"github.com/aarioai/airis/aa"
	"{{APP_BASE}}/cache"
	"{{APP_BASE}}/service"
	"sync"
	"time"
)

type Service struct {
	app   *core.App
	loc   *time.Location
	h     *cache.Cache
	mongo *mongodbhelper.Model
	s     *service.Service
}

var (
	svcOnce sync.Once
	svcObj  *Service
)

func New(app *core.App) *Service {
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
