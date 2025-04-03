package application

import (
	"github.com/follow1123/photos/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext interface {
	GetLogger() logger.AppLogger
	GetDB() *gorm.DB
}

type appContext struct {
	logger logger.AppLogger
	db     *gorm.DB
}

func NewAppContext(baseLogger *zap.SugaredLogger, db *gorm.DB) AppContext {
	return &appContext{
		logger: logger.NewAppLogger(baseLogger),
		db:     db,
	}
}

func (ac *appContext) GetLogger() logger.AppLogger {
	return ac.logger
}
func (ac *appContext) GetDB() *gorm.DB {
	return ac.db
}
