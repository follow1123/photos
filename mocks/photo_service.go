package mocks

import (
	"io"

	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/model/dto"
	"github.com/stretchr/testify/mock"
)

type PhotoService struct {
	mock.Mock
}

func (m *PhotoService) GetPhotoById(id uint) (*dto.PhotoDto, error) {
	ret := m.Called(id)

	var r0 *dto.PhotoDto
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*dto.PhotoDto)
	}

	r1 := ret.Error(1)
	return r0, r1
}

func (m *PhotoService) PhotoList() (*[]dto.PhotoDto, error) {
	ret := m.Called()

	var r0 *[]dto.PhotoDto
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]dto.PhotoDto)
	}

	r1 := ret.Error(1)
	return r0, r1
}

func (m *PhotoService) CreatePhoto(params []dto.CreatePhotoParam) error {
	ret := m.Called(params)
	return ret.Error(0)
}

func (m *PhotoService) UpdatePhoto(param dto.PhotoParam) (*dto.PhotoDto, error) {
	ret := m.Called(param)

	var r0 *dto.PhotoDto
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*dto.PhotoDto)
	}

	r1 := ret.Error(1)
	return r0, r1
}

func (m *PhotoService) DeletePhoto(id uint) error {
	ret := m.Called(id)
	return ret.Error(0)
}

func (m *PhotoService) GetPhotoFile(id uint, original bool) (io.ReadCloser, *imagemanager.ImageInfo, error) {
	ret := m.Called(id, original)

	var r0 io.ReadCloser
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(io.ReadCloser)
	}

	var r1 *imagemanager.ImageInfo
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*imagemanager.ImageInfo)
	}

	r2 := ret.Error(2)
	return r0, r1, r2
}
