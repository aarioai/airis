package core

import (
	"context"
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/config"
	"github.com/aarioai/airis/core/logger"
	"github.com/aarioai/airis/pkg/afmt"
)

type App struct {
	Config *config.Config
	Log    logger.LogInterface
}

func New(config *config.Config, logger logger.LogInterface) *App {
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
	if e.IsServerError() {
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
