package aa

import (
	"context"
	"github.com/aarioai/airis/aa/aconfig"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/aarioai/airis/pkg/afmt"
)

type App struct {
	Config        *aconfig.Config
	Log           alog.LogInterface
	LogWriter     *debug.LogWriter
	HTTP          *httpsvr.Settings
	GlobalContext acontext.Context
	GlobalCancel  context.CancelFunc
}

func New(ctx acontext.Context, cancel context.CancelFunc, config *aconfig.Config) *App {
	return &App{
		Config:        config,
		Log:           alog.NewDefaultLog(alog.LevelDebug),
		HTTP:          httpsvr.NewSettings(),
		GlobalContext: ctx,
		GlobalCancel:  cancel,
	}
}

func (app *App) WithLog(log alog.LogInterface) *App {
	app.Log = log
	return app
}

func (app *App) WithLogWriter(lw *debug.LogWriter) *App {
	app.LogWriter = lw
	return app
}

// Check checks errors.
func (app *App) Check(ctx context.Context, es ...*ae.Error) bool {
	e := afmt.First(es)
	if e == nil {
		return true
	}
	if e.IsServerError() || e.Detail != "" {
		app.Log.E(ctx, e)
	}
	return false
}

// CheckErrors checks standard errors.
func (app *App) CheckErrors(ctx context.Context, errs ...error) bool {
	err := afmt.First(errs)
	if err != nil && !ae.IsNotFound(err) {
		app.Log.Error(ctx, err.Error())
		return false
	}
	return true
}
