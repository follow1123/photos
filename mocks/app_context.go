package mocks

import (
	"github.com/follow1123/photos/config"
	imagemanager "github.com/follow1123/photos/imageManager"
	"github.com/follow1123/photos/logger"
	"github.com/stretchr/testify/mock"
)

type AppContext struct {
	mock.Mock
}

func (m *AppContext) GetLogger() *logger.AppLogger {
	ret := m.Called()
	var r0 *logger.AppLogger
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*logger.AppLogger)
	}
	return r0
}

func (m *AppContext) GetConfig() config.Config {
	ret := m.Called()

	var r0 config.Config
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(config.Config)
	}
	return r0
}

func (m *AppContext) GetImageManager() *imagemanager.ImageManager {
	ret := m.Called()

	var r0 *imagemanager.ImageManager
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*imagemanager.ImageManager)
	}
	return r0
}
