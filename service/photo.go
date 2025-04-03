package service

import (
	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/model"
	"github.com/follow1123/photos/model/dto"
	"gorm.io/gorm"
)

type PhotoService interface {
	GetPhotoById(uint) (*dto.PhotoDto, error)
	PhotoList() (*[]dto.PhotoDto, error)
	CreatePhoto(photo *dto.PhotoDto) (*dto.PhotoDto, error)
	UpdatePhoto(photo *dto.PhotoDto) (*dto.PhotoDto, error)
	DeletePhoto(uint) (*dto.PhotoDto, error)
}

type photoService struct {
	ctx application.AppContext
	db  *gorm.DB
}

func NewPhotoService(ctx application.AppContext) PhotoService {
	return &photoService{ctx: ctx, db: ctx.GetDB()}
}

func (serv *photoService) GetPhotoById(id uint) (*dto.PhotoDto, error) {
	var photo model.Photo
	if result := serv.db.First(&photo, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	photoDto := &dto.PhotoDto{}
	photoDto.Update(&photo)
	return photoDto, nil
}

func (serv *photoService) PhotoList() (*[]dto.PhotoDto, error) {
	var photoList []model.Photo

	if result := serv.db.Find(&photoList); result.Error != nil {
		return nil, result.Error
	}

	if len(photoList) == 0 {
		return nil, nil
	}

	dataList := make([]dto.PhotoDto, len(photoList))

	for i, v := range photoList {
		dataList[i] = dto.PhotoDto{}
		dataList[i].Update(&v)
	}

	return &dataList, nil
}

func (serv *photoService) CreatePhoto(photoDto *dto.PhotoDto) (*dto.PhotoDto, error) {
	p := photoDto.ToModel()
	if result := serv.db.Create(p); result.Error != nil {
		return nil, result.Error
	}
	photoDto.Update(p)
	return photoDto, nil
}

func (serv *photoService) UpdatePhoto(photoDto *dto.PhotoDto) (*dto.PhotoDto, error) {
	var photo model.Photo
	if result := serv.db.First(&photo, photoDto.ID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	photoDto.ID = 0
	if result := serv.db.Model(&photo).Updates(photoDto); result.Error != nil {
		return nil, result.Error
	}
	photoDto.Update(&photo)
	return photoDto, nil
}

func (serv *photoService) DeletePhoto(id uint) (*dto.PhotoDto, error) {
	var photo model.Photo
	if result := serv.db.First(&photo, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	if result := serv.db.Delete(&photo); result.Error != nil {
		return nil, result.Error
	}
	photoDto := &dto.PhotoDto{}
	photoDto.Update(&photo)
	return photoDto, nil
}
