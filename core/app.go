package core

import (
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
