package application

import (
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/filehandler"
	"github.com/follow1123/photos/logger"
	"go.uber.org/zap"
)

type AppContext interface {
	GetLogger() logger.AppLogger
	GetConfig() config.Config
	GetFileHandler(uri string) filehandler.FileHandler
}

type appContext struct {
	logger             logger.AppLogger
	config             config.Config
	fileHandlerFactory filehandler.FileHandlerFactory
}

func NewAppContext(baseLogger *zap.SugaredLogger, config config.Config) AppContext {
	appLogger := logger.NewAppLogger(baseLogger)
	fileHandlerFactory := filehandler.NewFileHandlerFactory(config.GetFilesPath(), appLogger)
	return &appContext{
		logger:             appLogger,
		config:             config,
		fileHandlerFactory: *fileHandlerFactory,
	}
}

func (c *appContext) GetLogger() logger.AppLogger {
	return c.logger
}

func (c *appContext) GetConfig() config.Config {
	return c.config
}

func (c *appContext) GetFileHandler(uri string) filehandler.FileHandler {
	return c.fileHandlerFactory.GetHandler(uri)
}
