package {{PACKAGE_NAME}}

import (
	"context"
	"github.com/aarioai/airis/core"
	"github.com/aarioai/airis/core/ae"
    "{{APP_BASE}}/cache"
    "{{APP_BASE}}/conf"
    "{{DRIVER_BASE}}/index"
    "{{DRIVER_BASE}}/mongodbhelper"
	"sync"
	"time"
)

type Service struct {
	app *core.App
	loc *time.Location
	h       *cache.Cache
	mongo *mongodbhelper.Model
}

var (
	svcOnce sync.Once
	svcObj  *Service

    mongoOnce sync.Once
    mongoObj  *mongodbhelper.Model
    mongoIndexEntities []index.Entity
)

func New(app *core.App) *Service {
	svcOnce.Do(func() {
		svcObj = &Service{
			app: app,
			loc: app.Config.TimeLocation,
            h:          cache.New(app),
            mongo:      NewMongo(app),
		}
	})
	return svcObj
}


func NewMongo(app *core.App) *mongodbhelper.Model {
	mongoOnce.Do(func() {
		mongoObj = mongodbhelper.NewDB(app, conf.MongoDBConfigSection)
		// create table and indexes
		if len(mongoIndexEntities) > 0 {
			for _, t := range mongoIndexEntities {
				ae.PanicOn(mongoObj.CreateIndexes(context.TODO(), t))
			}
		}
	})
	return mongoObj
}