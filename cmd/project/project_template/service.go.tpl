package service

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
	"{{APP_BASE}}/cache"
	"{{APP_BASE}}/conf"
	mbs "{{APP_BASE}}/module/bs/model"
	mcms "{{APP_BASE}}/module/cms/model"
	mss "{{APP_BASE}}/module/ss/model"
	"{{APP_BASE}}/private"
	"sync"
	"time"
)

type Service struct {
	app     *aa.App
	loc     *time.Location
	private *private.Service
	h       *cache.Cache
	mongo   *mongodb.Model
	mbs     *mbs.Model
	mcms    *mcms.Model
	mss     *mss.Model
}

var (
	svcOnce sync.Once
	svcObj  *Service
)

func New(app *aa.App) *Service {
	svcOnce.Do(func() {
		svcObj = &Service{
		    app:        app,
			loc:        app.Config.TimeLocation,
			private:    private.New(app),
			h:          cache.New(app),
			mongo:      mongodb.NewDB(app, conf.MongoCfgSection),
            mbs:        mbs.New(app),
            mcms:       mcms.New(app),
            mss:        mss.New(app),
		}
	})
	return svcObj
}

func (s *Service) mo(t index.Entity) *mongodb.ORMS {
	return s.mongo.ORM(t)
}

func (s *Service) Cache() *cache.Cache {
	return s.h
}

func (s *Service) Mongo() *mongodb.Model {
	return s.mongo
}