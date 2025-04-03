package migration

import (
	"github.com/follow1123/photos/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)
const VERSION = 0

const MIGRATION_PREFIX = "[MIGRATION] "

func InitOrMigration(db *gorm.DB, baseLogger *zap.SugaredLogger) {
	logger := baseLogger.Named(MIGRATION_PREFIX)
	logger.Debug("init or migration database")
	db.AutoMigrate(&model.Photo{})

}
