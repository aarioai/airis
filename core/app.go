package core

import (
	"context"
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/config"
	"github.com/aarioai/airis/core/logger"
)

type App struct {
	Config *config.Config
	Log    logger.LogInterface
}

func New(cfgPath string, logger logger.LogInterface) *App {
	c := config.New(cfgPath)
	return &App{
		Config: c,
		Log:    logger,
	}
}

// Check 检查错误
func (app *App) Check(ctx context.Context, es ...*ae.Error) bool {
	e := ae.Check(es...)
	if e == nil {
		return true
	}
	if e.IsServerError() {
		app.Log.Error(ctx, e.Text())
	}
	return false
}

// CheckError 检查标准错误
func (app *App) CheckErrors(ctx context.Context, errs ...error) bool {
	err := ae.CheckErrors(errs...)
	if err != nil {
		app.Log.Error(ctx, err.Error())
		return false
	}
	return true
}
