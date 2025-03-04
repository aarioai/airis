package model

import (
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis-driver/driver/sqlhelper"
	"github.com/aarioai/airis/aa"
	"{{APP_BASE}}/conf"
	"sync"
	"time"
)

type Model struct {
	app *core.App
	loc *time.Location
}

var (
	modelOnce          sync.Once
	modelObj           *Model
)

func New(app *core.App) *Model {
	modelOnce.Do(func() {
		modelObj = &Model{app: app, loc: app.Config.TimeLocation}
	})
	return modelObj
}

func (m *Model) db() *sqlhelper.DB {
	return sqlhelper.NewDriver(driver.NewMysqlPool(m.app, conf.MysqlConfigSection))
}


