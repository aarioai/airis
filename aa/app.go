package aa

import (
	"context"
	"github.com/aarioai/airis/aa/aconfig"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/pkg/afmt"
)

type App struct {
	Config        *aconfig.Config
	Log           alog.LogInterface
	GlobalContext acontext.Context
	GlobalCancel  context.CancelFunc
}

func New(ctx acontext.Context, cancel context.CancelFunc, config *aconfig.Config) *App {
	return &App{
		Config:        config,
		Log:           alog.NewDefaultLog(alog.LevelDebug),
		GlobalContext: ctx,
		GlobalCancel:  cancel,
	}
}

func (app *App) WithLog(log alog.LogInterface) *App {
	app.Log = log
	return app
}

// Check 检查错误
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

// CheckError 检查标准错误
func (app *App) CheckErrors(ctx context.Context, errs ...error) bool {
	err := afmt.First(errs)
	if err != nil && !ae.IsNotFound(err) {
		app.Log.Errorf(ctx, err.Error())
		return false
	}
	return true
}
