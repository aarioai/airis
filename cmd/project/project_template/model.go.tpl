package model

import (
	"github.com/aarioai/airis/core"
	"{{APP_BASE}}/conf"
	"{{DRIVER_BASE}}"
	"{{DRIVER_BASE}}/sqlhelper"
	"sync"
	"time"
)

type Model struct {
	app *core.App
	loc *time.Location
}

var (
	modelOnce sync.Once
	modelObj  *Model
)

func New(app *core.App) *Model {
	modelOnce.Do(func() {
		modelObj = &Model{app: app, loc: app.Config.TimeLocation}
	})
	return modelObj
}
func (m *Model) db() *sqlhelper.DB {
	return sqlhelper.AliveDriver(driver.NewMysql(m.app, conf.MysqlConfigSection))
}
