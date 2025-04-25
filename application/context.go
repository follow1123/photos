package application

import (
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
)

type AppContext struct {
	logger       *logger.AppLogger
	config       *config.Config
	imageManager *imagemanager.ImageManager
}

func NewAppContext(config *config.Config, imageManager *imagemanager.ImageManager, appLogger *logger.AppLogger) *AppContext {
	return &AppContext{
		logger:       appLogger,
		config:       config,
		imageManager: imageManager,
	}
}

func (ac *AppContext) GetLogger() *logger.AppLogger {
	return ac.logger
}

func (ac *AppContext) GetConfig() *config.Config {
	return ac.config
}

func (ac *AppContext) GetImageManager() *imagemanager.ImageManager {
	return ac.imageManager
}

func (ac *AppContext) Deinit() {
	ac.imageManager.Deinit()
}
