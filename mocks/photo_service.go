package mocks

import (
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

func (m *PhotoService) CreatePhoto(photoDto *dto.PhotoDto) (*dto.PhotoDto, error) {
	ret := m.Called(photoDto)

	var r0 *dto.PhotoDto
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*dto.PhotoDto)
	}

	r1 := ret.Error(1)
	return r0, r1
}

func (m *PhotoService) UpdatePhoto(photoDto *dto.PhotoDto) (*dto.PhotoDto, error) {
	ret := m.Called(photoDto)

	var r0 *dto.PhotoDto
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*dto.PhotoDto)
	}

	r1 := ret.Error(1)
	return r0, r1
}

func (m *PhotoService) DeletePhoto(id uint) (*dto.PhotoDto, error) {
	ret := m.Called(id)

	var r0 *dto.PhotoDto
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*dto.PhotoDto)
	}

	r1 := ret.Error(1)
	return r0, r1
}
