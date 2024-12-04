package core

import (
	"context"
	"github.com/aarioai/airis/core/ae"
)

// 快捷方式，对服务器错误记录日志
func (app *App) Check(ctx context.Context, e *ae.Error) bool {
	if e == nil {
		return true
	}
	if e.IsServerError() {
		app.Log.Error(ctx, e.Text())
	}
	return false
}

func (app *App) CheckError(ctx context.Context, err error) bool {
	if err == nil {
		return true
	}
	app.Log.Error(ctx, err.Error())
	return false
}

// 快捷panic
func (app *App) Assert(ctx context.Context, e *ae.Error) {
	if e != nil && e.IsServerError() {
		app.Log.Error(ctx, e.Text())
		panic(e.Text())
	}
}

func (app *App) AssertError(ctx context.Context, err error) {
	if err != nil {
		app.Log.Error(ctx, err.Error())
		panic(err.Error())
	}
}
