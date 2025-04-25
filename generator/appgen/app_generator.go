package appgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/database"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/webserver"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AppComponents struct {
	BaseLogger   *zap.SugaredLogger
	GinLogger    *logger.GinLogger
	GormLogger   *logger.GormLogger
	AppLogger    *logger.AppLogger
	Config       *config.Config
	DB           *database.SqliteDB
	ImageCache   *imagemanager.ImageCache
	ImageManager *imagemanager.ImageManager
	AppContext   *application.AppContext
	WebServer    *webserver.GinWebServer
}

func GenConfig(appComponents *AppComponents) (*config.Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var path = filepath.Join(wd, fmt.Sprintf("test_%s", strings.ReplaceAll(uuid.New().String(), "-", "")))
	conf := config.NewConfig(
		config.WithAddress(":8088"),
		config.WithPath(path),
	)
	appComponents.Config = conf
	return conf, nil
}

func GenBaseLogger(appComponents *AppComponents) (*zap.SugaredLogger, error) {
	baseLogger, err := logger.NewBaseLogger()
	if err != nil {
		return nil, err
	}
	appComponents.BaseLogger = baseLogger
	return baseLogger, nil
}

func GenGinLogger(appComponents *AppComponents) (*logger.GinLogger, error) {
	if appComponents.BaseLogger == nil {
		baseLogger, err := GenBaseLogger(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.BaseLogger = baseLogger
	}

	return logger.NewGinLogger(appComponents.BaseLogger), nil
}

func GenAppLogger(appComponents *AppComponents) (*logger.AppLogger, error) {
	if appComponents.BaseLogger == nil {
		baseLogger, err := GenBaseLogger(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.BaseLogger = baseLogger
	}

	return logger.NewAppLogger(appComponents.BaseLogger), nil
}

func GenGormLogger(appComponents *AppComponents) (*logger.GormLogger, error) {
	if appComponents.BaseLogger == nil {
		baseLogger, err := GenBaseLogger(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.BaseLogger = baseLogger
	}

	return logger.NewGormLogger(appComponents.BaseLogger), nil
}

func GenDatabase(appComponents *AppComponents) (*database.SqliteDB, error) {
	if appComponents.Config == nil {
		config, err := GenConfig(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.Config = config
	}

	if appComponents.GormLogger == nil {
		gormLogger, err := GenGormLogger(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.GormLogger = gormLogger
	}

	err := appComponents.Config.CreatePath()
	if err != nil {
		return nil, err
	}

	db, err := database.NewDatabase(appComponents.Config, appComponents.GormLogger)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GenDBMigrator(appComponents *AppComponents) (*database.DBMigrator, error) {
	if appComponents.DB == nil {
		db, err := GenDatabase(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.DB = db
	}

	return database.NewDBMigrator(appComponents.DB), nil
}

func GenWebServer(appComponents *AppComponents) (*webserver.GinWebServer, error) {
	if appComponents.Config == nil {
		config, err := GenConfig(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.Config = config
	}

	if appComponents.GinLogger == nil {
		ginLogger, err := GenGinLogger(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.GinLogger = ginLogger
	}

	return webserver.NewGinWebServer(appComponents.Config, appComponents.GinLogger), nil
}

func GenImageCache(appComponents *AppComponents) (*imagemanager.ImageCache, error) {
	return imagemanager.NewImageCache()
}

func GenImageManager(appComponents *AppComponents) (*imagemanager.ImageManager, error) {
	if appComponents.Config == nil {
		config, err := GenConfig(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.Config = config
	}

	if appComponents.ImageCache == nil {
		imageCache, err := GenImageCache(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.ImageCache = imageCache
	}

	if appComponents.AppLogger == nil {
		appLogger, err := GenAppLogger(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.AppLogger = appLogger
	}

	return imagemanager.NewImageManager(appComponents.Config.GetFilesPath(), appComponents.ImageCache, appComponents.AppLogger), nil
}

func GenAppContext(appComponents *AppComponents) (*application.AppContext, error) {
	if appComponents.Config == nil {
		config, err := GenConfig(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.Config = config
	}

	if appComponents.ImageManager == nil {
		imageMgr, err := GenImageManager(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.ImageManager = imageMgr
	}

	if appComponents.AppLogger == nil {
		appLogger, err := GenAppLogger(appComponents)
		if err != nil {
			return nil, err
		}
		appComponents.AppLogger = appLogger
	}

	return application.NewAppContext(appComponents.Config, appComponents.ImageManager, appComponents.AppLogger), nil
}
