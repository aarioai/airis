package {{APP_NAME}}

import (
	"github.com/aarioai/airis/aa"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type Service struct {
	app    *aa.App
	loc    *time.Location
	mtx     sync.RWMutex
	conn   *grpc.ClientConn
	target string
}

var (
	svcOnce sync.Once
	svcObj  *Service
)

func New(app *aa.App) *Service {
	svcOnce.Do(func() {
		svcObj = &Service{
			app: app,
			loc: app.Config.TimeLocation,
		}

	})
	return svcObj
}
