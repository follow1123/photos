package database

import (
	"github.com/follow1123/photos/model"
)

const VERSION = 0

type DBMigrator struct {
	DB *SqliteDB
}

func NewDBMigrator(db *SqliteDB) *DBMigrator {
	return &DBMigrator{DB: db}
}

func (dm *DBMigrator) InitOrMigrate() {
	dm.DB.Logger.Logger.Info("DATABASE MIGRATION START")
	defer dm.DB.Logger.Logger.Info("DATABASE MIGRATION END")
	dm.DB.AutoMigrate(&model.Photo{})
}
