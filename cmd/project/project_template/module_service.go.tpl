package {{MODULE_NAME}}

import (
	"github.com/aarioai/airis/core"
	"{{APP_BASE}}/cache"
	"{{APP_BASE}}/module/{{MODULE_NAME}}/model"
	"{{APP_BASE}}/service"
	"sync"
	"time"
)

type Service struct {
    app *core.App
	loc *time.Location
	h   *cache.Cache
	m   *model.Model
	s   *service.Service
}

var (
	svcOnce sync.Once
	svcObj  *Service
)

func New(app *core.App) *Service {
	svcOnce.Do(func() {
		svcObj = &Service{
			app: app,
			loc: app.Config.TimeLocation,
			h:   cache.New(app),
			m:   model.New(app),
			s:   service.New(app),
		}
	})
	return svcObj
}
