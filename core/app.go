package core

import (
	"context"
	"github.com/aarioai/airis/core/aconfig"
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/alog"
	"github.com/aarioai/airis/pkg/afmt"
)

type App struct {
	Config *aconfig.Config
	Log    alog.LogInterface
}

func New(config *aconfig.Config, logger alog.LogInterface) *App {
	return &App{
		Config: config,
		Log:    logger,
	}
}

// Check 检查错误
func (app *App) Check(ctx context.Context, es ...*ae.Error) bool {
	e := afmt.First(es)
	if e == nil {
		return true
	}
	if e.IsServerError() || e.Detail != "" {
		app.Log.Error(ctx, e.Text())
	}
	return false
}

// CheckError 检查标准错误
func (app *App) CheckErrors(ctx context.Context, errs ...error) bool {
	err := afmt.First(errs)
	if err != nil {
		app.Log.Error(ctx, err.Error())
		return false
	}
	return true
}
