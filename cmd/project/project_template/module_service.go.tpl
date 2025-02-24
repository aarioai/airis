package {{MODULE_NAME}}

import (
	"github.com/aarioai/airis/core"
	"{{APP_BASE}}/cache"
	"{{APP_BASE}}/module/{{MODULE_NAME}}/model"
	"{{APP_BASE}}/service"
	"{{DRIVER_BASE}}/mongodbhelper"
	"sync"
	"time"
)

type Service struct {
	app   *core.App
	loc   *time.Location
	h     *cache.Cache
	mongo *mongodbhelper.Model
	m     *model.Model
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
			m:     model.New(app),
			s:     s,
		}
	})
	return svcObj
}
