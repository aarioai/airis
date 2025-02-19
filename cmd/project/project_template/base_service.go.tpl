package {{PACKAGE_NAME}}

import (
	"github.com/aarioai/airis/core"
	"sync"
	"time"
)

type Service struct {
	app *core.App
	loc *time.Location
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
		}
	})
	return svcObj
}
