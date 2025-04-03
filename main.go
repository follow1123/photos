package main

import (
	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/middleware"
	"github.com/follow1123/photos/migration"
	"github.com/follow1123/photos/router"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	baseLogger := logger.InitBaseLogger()
	db, err := gorm.Open(sqlite.Open("photos.db"), &gorm.Config{
		Logger: logger.NewGormLogger(baseLogger),
	})
	if err != nil {
		baseLogger.Fatal("cannot connect to sqlite database")
	}

	migration.InitOrMigration(db, baseLogger)

	appCtx := application.NewAppContext(baseLogger, db)

	r := gin.New()
	// middleware
	r.Use(logger.NewGinLoggerHandler(baseLogger), gin.Recovery())
	r.Use(middleware.GlobalErrorHandler(baseLogger))

	// router
	router.Init(r, appCtx, baseLogger)

	r.Run()
}
