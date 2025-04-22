package application

import (
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"go.uber.org/zap"
)

type AppContext interface {
	GetLogger() *logger.AppLogger
	GetConfig() config.Config
	GetImageManager() *imagemanager.ImageManager
	Deinit()
}

type appContext struct {
	logger       *logger.AppLogger
	config       config.Config
	imageManager *imagemanager.ImageManager
}

func NewAppContext(baseLogger *zap.SugaredLogger, config config.Config) AppContext {
	appLogger := logger.NewAppLogger(baseLogger)
	imageManager := imagemanager.NewImageManager(config.GetFilesPath(), appLogger)
	return &appContext{
		logger:       appLogger,
		config:       config,
		imageManager: imageManager,
	}
}

func (c *appContext) GetLogger() *logger.AppLogger {
	return c.logger
}

func (c *appContext) GetConfig() config.Config {
	return c.config
}

func (c *appContext) GetImageManager() *imagemanager.ImageManager {
	return c.imageManager
}

func (c *appContext) Deinit() {
	c.imageManager.Deinit()
}
