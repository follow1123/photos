package application

import (
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/logger"
	"go.uber.org/zap"
)

type AppContext interface {
	GetLogger() *logger.AppLogger
	GetConfig() config.Config
}

type appContext struct {
	logger             *logger.AppLogger
	config             config.Config
}

func NewAppContext(baseLogger *zap.SugaredLogger, config config.Config) AppContext {
	appLogger := logger.NewAppLogger(baseLogger)
	return &appContext{
		logger:             appLogger,
		config:             config,
	}
}

func (c *appContext) GetLogger() *logger.AppLogger {
	return c.logger
}

func (c *appContext) GetConfig() config.Config {
	return c.config
}
