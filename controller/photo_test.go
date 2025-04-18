package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/middleware"
	"github.com/follow1123/photos/mocks"
	"github.com/follow1123/photos/model/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PhotoAPISuite struct {
	suite.Suite
	r    *gin.Engine
	serv *mocks.PhotoService
	ctl  *photoController
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &PhotoAPISuite{})
}

func (s *PhotoAPISuite) SetupSuite() {
	r := gin.Default()
	zapLogger := zap.NewExample().Sugar()
	r.Use(middleware.GlobalErrorHandler(zapLogger))

	cfg := config.NewConfig(":8080")
	appCtx := application.NewAppContext(zapLogger, cfg)

	photoCtl := NewPhotoController(appCtx, nil)
	r.GET(PHOTO_API_GETBYID, photoCtl.GetPhotoById)
	r.GET(PHOTO_API_LIST, photoCtl.PhotoList)
	r.POST(PHOTO_API_CREATE, photoCtl.CreatePhoto)
	r.PUT(PHOTO_API_UPDATE, photoCtl.UpdatePhoto)
	r.DELETE(PHOTO_API_DELETE, photoCtl.DeletePhoto)
	s.ctl, _ = photoCtl.(*photoController)
	s.r = r
}

func (s *PhotoAPISuite) SetupTest() {
	service := mocks.PhotoService{}
	s.serv = &service
	s.ctl.service = &service
}

func (s *PhotoAPISuite) TestGetByIdSuccess() {
	expectedCode, expectedData := http.StatusOK, dto.PhotoDto{
		Desc: "2343214",
		Uri:  "123412",
	}
	s.serv.On("GetPhotoById", mock.Anything).Return(&expectedData, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/photo/1", nil)
	s.r.ServeHTTP(w, req)

	s.Equal(expectedCode, w.Code)
	expectedDataJson, _ := json.Marshal(expectedData)
	s.Equal(string(expectedDataJson), w.Body.String())
}

func (s *PhotoAPISuite) TestGetByIdFailure() {
	scenarios := []struct {
		uri          string
		err          error
		expectedCode int
		expectedBody string
	}{
		{"/photo/a", nil, 400, `{"message":"strconv.ParseUint: parsing \"a\": invalid syntax"}`},
		{"/photo/1", nil, 404, `{"message":"data not found"}`},
		{"/photo/1", gorm.ErrDuplicatedKey, 500, `{"message":"Internal Server Error"}`},
	}

	for _, scenario := range scenarios {
		w := httptest.NewRecorder()
		s.serv.On("GetPhotoById", mock.Anything).Return(nil, scenario.err)
		req, _ := http.NewRequest("GET", scenario.uri, nil)
		s.r.ServeHTTP(w, req)
		s.Equal(scenario.expectedCode, w.Code)
		s.Equal(scenario.expectedBody, w.Body.String())
		s.serv.On("GetPhotoById").Unset()
	}
}

func (s *PhotoAPISuite) TestDeletePhotoSuccess() {
	expectedCode, expectedData := http.StatusOK, dto.PhotoDto{
		Desc: "2343214",
		Uri:  "123412",
	}
	s.serv.On("DeletePhoto", mock.Anything).Return(&expectedData, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/photo/1", nil)
	s.r.ServeHTTP(w, req)

	s.Equal(expectedCode, w.Code)

	expectedDataJson, _ := json.Marshal(expectedData)
	s.Equal(string(expectedDataJson), w.Body.String())
}

func (s *PhotoAPISuite) TestDeletePhotoFailure() {
	scenarios := []struct {
		uri          string
		err          error
		expectedCode int
		expectedBody string
	}{
		{"/photo/a", nil, 400, `{"message":"strconv.ParseUint: parsing \"a\": invalid syntax"}`},
		{"/photo/1", nil, 404, `{"message":"data not found"}`},
		{"/photo/1", gorm.ErrDuplicatedKey, 500, `{"message":"Internal Server Error"}`},
	}

	for _, scenario := range scenarios {
		w := httptest.NewRecorder()
		s.serv.On("DeletePhoto", mock.Anything).Return(nil, scenario.err)
		req, _ := http.NewRequest("DELETE", scenario.uri, nil)
		s.r.ServeHTTP(w, req)
		s.Equal(scenario.expectedCode, w.Code)
		s.Equal(scenario.expectedBody, w.Body.String())
		s.serv.On("DeletePhoto").Unset()
	}
}
