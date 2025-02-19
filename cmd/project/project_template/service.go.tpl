package service

import (
	"github.com/aarioai/airis/core"
	"{{APP_BASE}}/cache"
	mbs "{{APP_BASE}}/module/bs/model"
	mcns "{{APP_BASE}}/module/cns/model"
	mss "{{APP_BASE}}/module/ss/model"
	"{{APP_BASE}}/private"
	"sync"
	"time"
)

type Service struct {
	app     *core.App
	loc     *time.Location
	private *private.Service
	h       *cache.Cache
	mbs      *mbs.Model
	mcns     *mcns.Model
	mss      *mss.Model
}

var (
	svcOnce sync.Once
	svcObj  *Service
)

func New(app *core.App) *Service {
	svcOnce.Do(func() {
		svcObj = &Service{
		    app:     app,
			loc:     app.Config.TimeLocation,
			private: private.New(app),
			h:       cache.New(app),
			bs:      bs.New(app),
			cns:     cns.New(app),
			ss:      ss.New(app),
            bsm:      bsm.New(app),
            mcns:     mcns.New(app),
            mss:      mss.New(app),
		}
	})
	return svcObj
}
