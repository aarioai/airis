package {{PACKAGE_NAME}}

import (
    "github.com/aarioai/airis-driver/driver/mongodbhelper"
	"github.com/aarioai/airis/core"
    "{{APP_BASE}}/cache"
    "{{APP_BASE}}/conf"
	"sync"
	"time"
)

type Service struct {
	app   *core.App
	loc   *time.Location
	h     *cache.Cache
	mongo *mongodbhelper.Model
}

var (
	svcOnce sync.Once
	svcObj  *Service
)

func New(app *core.App) *Service {
	svcOnce.Do(func() {
		svcObj = &Service{
			app:   app,
			loc:   app.Config.TimeLocation,
            h:     cache.New(app),
            mongo: mongodbhelper.NewDB(app, conf.MongoDBConfigSection),
		}
	})
	return svcObj
}