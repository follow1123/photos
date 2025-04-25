package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/controller"
	"github.com/follow1123/photos/generator/appgen"
	"github.com/follow1123/photos/mocks"
	"github.com/follow1123/photos/model/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PhotoAPISuite struct {
	suite.Suite
	r    *gin.Engine
	serv *mocks.PhotoService
}

func TestPhotoAPISuite(t *testing.T) {
	suite.Run(t, &PhotoAPISuite{})
}

func (s *PhotoAPISuite) SetupSuite() {
	appComponents := &appgen.AppComponents{}
	ctx, err := appgen.GenAppContext(appComponents)
	s.Nil(err)
	ws, err := appgen.GenWebServer(appComponents)
	s.Nil(err)
	s.serv = &mocks.PhotoService{}

	ws.InitMiddleware()

	ws.SetRouters(
		controller.NewPhotoController(ctx, s.serv),
	)
	ws.InitRouter()

	s.r = ws.GetEngine()
}

func (s *PhotoAPISuite) TestGetByIdSuccess() {
	expectedCode, expectedData := http.StatusOK, dto.PhotoDto{
		ID:   1,
		Desc: "2343214",
	}
	s.serv.On("GetPhotoById", mock.Anything).Return(&expectedData, nil)
	defer s.serv.On("GetPhotoById").Unset()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/photo/1", nil)
	s.r.ServeHTTP(w, req)

	s.Equal(expectedCode, w.Code)
	expectedDataJson, err := json.Marshal(expectedData)
	s.Nil(err)
	s.Equal(string(expectedDataJson), w.Body.String())
}

func (s *PhotoAPISuite) TestGetByIdFailure() {
	_, convErr := strconv.ParseUint("a", 10, 32)
	scenarios := []struct {
		uri          string
		err          error
		expectedCode int
		expectedData application.AppError
	}{
		{"/photo/a", nil, 400, application.AppError{Message: convErr.Error()}},
		{"/photo/1", application.ErrDataNotFound, 404, *application.ErrDataNotFound},
		{"/photo/1", gorm.ErrDuplicatedKey, 500, *application.ErrInternalServerError},
	}

	for _, scenario := range scenarios {
		w := httptest.NewRecorder()
		s.serv.On("GetPhotoById", mock.Anything).Return(nil, scenario.err)
		req, _ := http.NewRequest("GET", scenario.uri, nil)
		s.r.ServeHTTP(w, req)
		s.Equal(scenario.expectedCode, w.Code)

		expectedBody, err := json.Marshal(scenario.expectedData)
		s.Nil(err)
		s.Equal(string(expectedBody), w.Body.String())
		s.serv.On("GetPhotoById").Unset()
	}
}

func (s *PhotoAPISuite) TestDeletePhotoSuccess() {
	expectedCode := http.StatusNoContent
	s.serv.On("DeletePhoto", mock.Anything).Return(nil)
	defer s.serv.On("DeletePhoto").Unset()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/photo/1", nil)
	s.r.ServeHTTP(w, req)

	s.Equal(expectedCode, w.Code)

	s.Equal(0, w.Body.Len())
}

func (s *PhotoAPISuite) TestDeletePhotoFailure() {
	_, convErr := strconv.ParseUint("a", 10, 32)
	scenarios := []struct {
		uri          string
		err          error
		expectedCode int
		expectedData application.AppError
	}{
		{"/photo/a", nil, 400, application.AppError{Message: convErr.Error()}},
		{"/photo/1", application.ErrDataNotFound, 404, *application.ErrDataNotFound},
		{"/photo/1", gorm.ErrDuplicatedKey, 500, *application.ErrInternalServerError},
	}

	for _, scenario := range scenarios {
		w := httptest.NewRecorder()
		s.serv.On("DeletePhoto", mock.Anything).Return(scenario.err)
		req, _ := http.NewRequest("DELETE", scenario.uri, nil)
		s.r.ServeHTTP(w, req)
		s.Equal(scenario.expectedCode, w.Code)

		expectedBody, err := json.Marshal(scenario.expectedData)
		s.Nil(err)
		s.Equal(string(expectedBody), w.Body.String())
		s.serv.On("DeletePhoto").Unset()
	}
}
