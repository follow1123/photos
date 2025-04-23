package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/model"
	"github.com/follow1123/photos/model/dto"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PhotoServiceSuite struct {
	suite.Suite
	serv   PhotoService
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func TestPhotoServiceSuite(t *testing.T) {
	suite.Run(t, &PhotoServiceSuite{})
}

func (s *PhotoServiceSuite) SetupSuite() {
	baseLogger := logger.InitBaseLogger()
	s.logger = baseLogger.Named("PHOTO_SERVICE_TEST")
	db, err := gorm.Open(sqlite.Open("photos_test.db"), &gorm.Config{
		Logger: logger.NewGormLogger(baseLogger),
	})
	if err != nil {
		baseLogger.Fatal("cannot connect to sqlite database")
	}

	cfg := config.NewConfig("")

	zapLogger := zap.NewExample().Sugar()
	appCtx := application.NewAppContext(zapLogger, cfg)
	serv := NewPhotoService(appCtx, db)
	s.serv = serv
	s.db = db
}

func (s *PhotoServiceSuite) TearDownSuite() {
	session, _ := s.db.DB()
	session.Close()
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	os.Remove(filepath.Join(dir, "photos_test.db"))
}

func (s *PhotoServiceSuite) SetupTest() {
	s.db.Migrator().CreateTable(&model.Photo{})
}

func (s *PhotoServiceSuite) TearDownTest() {
	s.db.Migrator().DropTable(&model.Photo{})
}

func (s *PhotoServiceSuite) TestGetByIdSuccess() {
	expectedData := &dto.PhotoDto{
		ID:        1,
		Desc:      "2343214",
		Size:      102400,
		PhotoDate: time.Now(),
	}
	expectedDataJson, _ := json.Marshal(expectedData)
	s.db.Create(expectedData.ToModel())
	data, err := s.serv.GetPhotoById(1)
	dataJson, _ := json.Marshal(data)
	s.Nil(err)
	s.Equal(string(expectedDataJson), string(dataJson))

}

func (s *PhotoServiceSuite) TestGetByIdFailure() {
	expectedErr := application.ErrDataNotFound
	data, err := s.serv.GetPhotoById(1)
	s.Nil(data)
	s.Equal(expectedErr, err)
}

func (s *PhotoServiceSuite) TestDeletePhotoSuccess() {
	expectedData := &dto.PhotoDto{
		ID:        1,
		Desc:      "2343214",
		Size:      102400,
		PhotoDate: time.Now(),
	}
	s.db.Create(expectedData.ToModel())

	err := s.serv.DeletePhoto(1)
	s.Nil(err)
}

func (s *PhotoServiceSuite) TestDeletePhotoFailure() {
	err := s.serv.DeletePhoto(1)
	s.NotNil(err)
}
