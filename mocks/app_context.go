package mocks

import (
	"github.com/follow1123/photos/logger"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type AppContext struct {
	mock.Mock
}

func (m *AppContext) GetLogger() logger.AppLogger {
	ret := m.Called()
	var r0 logger.AppLogger
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(logger.AppLogger)
	}
	return r0
}
func (m *AppContext) GetDB() *gorm.DB {
	ret := m.Called()

	var r0 *gorm.DB
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*gorm.DB)
	}
	return r0
}
