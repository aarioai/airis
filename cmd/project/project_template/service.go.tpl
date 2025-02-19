package service

import (
	"context"
	"github.com/aarioai/airis/core"
	"github.com/aarioai/airis/core/ae"
	"{{APP_BASE}}/cache"
	"{{APP_BASE}}/conf"
	mbs "{{APP_BASE}}/module/bs/model"
	mcns "{{APP_BASE}}/module/cns/model"
	mss "{{APP_BASE}}/module/ss/model"
	"{{APP_BASE}}/private"
    "{{DRIVER_BASE}}/index"
    "{{DRIVER_BASE}}/mongodbhelper"
	"sync"
	"time"
)

type Service struct {
	app     *core.App
	loc     *time.Location
	private *private.Service
	h       *cache.Cache
	mongo *mongodbhelper.Model
	mbs      *mbs.Model
	mcns     *mcns.Model
	mss      *mss.Model
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
		    app:        app,
			loc:        app.Config.TimeLocation,
			private:    private.New(app),
			h:          cache.New(app),
			mongo:      NewMongo(app),
            mbs:        mbs.New(app),
            mcns:       mcns.New(app),
            mss:        mss.New(app),
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