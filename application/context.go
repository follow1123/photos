package application

import (
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"go.uber.org/zap"
)

type AppContext struct {
	logger       *logger.AppLogger
	config       *config.Config
	imageManager *imagemanager.ImageManager
}

func NewAppContext(baseLogger *zap.SugaredLogger, config *config.Config) *AppContext {
	appLogger := logger.NewAppLogger(baseLogger)
	imageManager := imagemanager.NewImageManager(config.GetFilesPath(), appLogger)
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
