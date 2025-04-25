package main

import (
	"fmt"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/controller"
	"github.com/follow1123/photos/database"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/service"
	"github.com/follow1123/photos/webserver"
)

func main() {
	// logger
	baseLogger, err := logger.NewBaseLogger()
	if err != nil {
		panic(fmt.Sprintf("init base logger error: %v", err))
	}

	appLogger := logger.NewAppLogger(baseLogger)
	ginLogger := logger.NewGinLogger(baseLogger)
	gormLogger := logger.NewGormLogger(baseLogger)

	// config
	conf := config.NewConfig()
	err = conf.CreatePath()
	if err != nil {
		panic(fmt.Sprintf("cannot create config path: %s, error: %v", conf.GetPrefixPath(), err))
	}

	// database
	db, err := database.NewDatabase(conf, gormLogger)
	if err != nil {
		panic(fmt.Sprintf("init database error: %v", err))
	}

	// migration
	dbMigrator := database.NewDBMigrator(db)
	dbMigrator.InitOrMigrate()

	// image manager
	imageCache, err := imagemanager.NewImageCache()
	if err != nil {
		panic(fmt.Sprintf("init image cache error: %v", err))
	}

	imageManager := imagemanager.NewImageManager(conf.GetFilesPath(), imageCache, appLogger)

	// application context
	appCtx := application.NewAppContext(conf, imageManager, appLogger)
	defer appCtx.Deinit()

	// web server
	ws := webserver.NewGinWebServer(conf, ginLogger)

	// middleware
	ws.InitMiddleware()

	// router
	photoServ := service.NewPhotoService(appCtx, db)

	ws.SetRouters(
		controller.NewPhotoController(appCtx, photoServ),
	)

	ws.InitRouter()

	ws.Start()
}
