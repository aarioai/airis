package {{MODULE_NAME}}

import (
	"github.com/aarioai/airis-driver/driver/mongodbhelper"
	"github.com/aarioai/airis/aa"
	"{{APP_BASE}}/cache"
	"{{APP_BASE}}/module/{{MODULE_NAME}}/model"
	"{{APP_BASE}}/service"
	"sync"
	"time"
)

type Service struct {
	app   *aa.App
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

func New(app *aa.App) *Service {
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
