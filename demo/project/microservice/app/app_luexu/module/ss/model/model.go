package model

import (
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis-driver/driver/mysqli"
	"github.com/aarioai/airis/aa"
	"project/microservice/app/app_luexu/conf"
	"sync"
	"time"
)

type Model struct {
	app *aa.App
	loc *time.Location
}

var (
	modelOnce          sync.Once
	modelObj           *Model
)

func New(app *aa.App) *Model {
	modelOnce.Do(func() {
		modelObj = &Model{app: app, loc: app.Config.TimeLocation}
	})
	return modelObj
}

func (m *Model) db() *mysqli.DB {
	return mysqli.NewDriver(driver.NewMysqlPool(m.app, conf.MysqlCfgSection))
}


