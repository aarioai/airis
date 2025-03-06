package service

import (
    "github.com/aarioai/airis-driver/driver/mongodbhelper"
	"github.com/aarioai/airis/aa"
	"{{APP_BASE}}/cache"
	"{{APP_BASE}}/conf"
	mbs "{{APP_BASE}}/module/bs/model"
	mcns "{{APP_BASE}}/module/cns/model"
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
	mongo   *mongodbhelper.Model
	mbs     *mbs.Model
	mcns    *mcns.Model
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
			mongo:      mongodbhelper.NewDB(app, conf.MongoDBConfigSection),
            mbs:        mbs.New(app),
            mcns:       mcns.New(app),
            mss:        mss.New(app),
		}
	})
	return svcObj
}

func (s *Service) Cache() *cache.Cache {
	return s.h
}

func (s *Service) Mongo() *mongodbhelper.Model {
	return s.mongo
}

func (s *Service) Private() *private.Service {
	return s.private
}