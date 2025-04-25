package database

import (
	"os"
	"path/filepath"

	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDB struct {
	*gorm.DB
	Config *config.Config
	DBFile string
	Logger *logger.GormLogger
}

func NewDatabase(config *config.Config, gormLogger *logger.GormLogger) (*SqliteDB, error) {
	dbFile := filepath.Join(config.GetPrefixPath(), "photos.db")
	db, err := gorm.Open(
		sqlite.Open(dbFile),
		&gorm.Config{Logger: gormLogger},
	)
	if err != nil {
		return nil, err
	}

	return &SqliteDB{DB: db, Config: config, DBFile: dbFile, Logger: gormLogger}, nil
}

func (d *SqliteDB) DeleteDBFile() error {
	err := os.Remove(d.DBFile)
	if err != nil {
		return err
	}
	return nil
}
