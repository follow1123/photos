package service_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/database"
	"github.com/follow1123/photos/generator/appgen"
	"github.com/follow1123/photos/generator/imagegen"
	"github.com/follow1123/photos/generator/valuegen"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/model"
	"github.com/follow1123/photos/model/dto"
	"github.com/follow1123/photos/service"
	"github.com/stretchr/testify/suite"
)

type PhotoServiceSuite struct {
	suite.Suite
	serv   service.PhotoService
	db     *database.SqliteDB
	config *config.Config
}

func TestPhotoServiceSuite(t *testing.T) {
	suite.Run(t, &PhotoServiceSuite{})
}

func (s *PhotoServiceSuite) SetupSuite() {
	appComponents := &appgen.AppComponents{}
	ctx, err := appgen.GenAppContext(appComponents)
	s.Nil(err)
	db, err := appgen.GenDatabase(appComponents)
	s.Nil(err)

	migrator, err := appgen.GenDBMigrator(appComponents)
	s.Nil(err)
	migrator.InitOrMigrate()

	serv := service.NewPhotoService(ctx, db)

	s.serv = serv
	s.db = db
	s.config = appComponents.Config
}

func (s *PhotoServiceSuite) TearDownSuite() {
	session, err := s.db.DB.DB()
	s.Nil(err)
	session.Close()
	s.config.DeletePath()
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

func (s *PhotoServiceSuite) TestCreatePhotoSuccess() {
	var recordCount = 3
	params := make([]dto.CreatePhotoParam, 0, recordCount)
	for i := range recordCount {
		_ = i

		var param = dto.CreatePhotoParam{
			Desc:      valuegen.GenString(valuegen.WithStrLimit(5, 30)),
			PhotoDate: *valuegen.GenTime(nil, nil),
		}

		buf := new(bytes.Buffer)
		_, err := imagegen.GenImage(buf)
		s.Nil(err)
		param.ImageSource = imagemanager.NewReaderSource(bytes.NewReader(buf.Bytes()), param.Desc)
		params = append(params, param)
	}
	failureResults := s.serv.CreatePhoto(params)
	s.True(len(failureResults) == 0)
}

func (s *PhotoServiceSuite) TestCreatePhotoFailure() {
	var i uint = 1

	buildImage := func() []byte {
		buf := new(bytes.Buffer)
		_, err := imagegen.GenImage(buf)
		s.Nil(err)
		return buf.Bytes()
	}

	buildParam := func() dto.CreatePhotoParam {
		i++
		return dto.CreatePhotoParam{
			UploadID:  i,
			Desc:      valuegen.GenString(valuegen.WithStrLimit(5, 30)),
			PhotoDate: *valuegen.GenTime(nil, nil),
		}
	}

	name1 := "aaa"
	name2 := "bbb"
	img1 := buildImage()

	var param1 = buildParam()
	param1.ImageSource = imagemanager.NewReaderSource(bytes.NewReader(img1), name1)
	var param2 = buildParam()
	param2.ImageSource = imagemanager.NewReaderSource(bytes.NewReader(img1), name2)

	var params1 = []dto.CreatePhotoParam{param1, param2}

	scenarios := []struct {
		params       []dto.CreatePhotoParam
		checkResults func([]dto.CreatePhotoFailedResult)
	}{
		{params1, func(results []dto.CreatePhotoFailedResult) {
			s.True(len(results) == 1)
			r := results[0]
			s.True(strings.Contains(r.Message, "重复"))
		}},
	}
	for _, scenario := range scenarios {
		failureResults := s.serv.CreatePhoto(scenario.params)
		scenario.checkResults(failureResults)
	}

}
