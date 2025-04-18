package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/model"
	"github.com/follow1123/photos/model/dto"
	"gorm.io/gorm"
)

type PhotoService interface {
	GetPhotoById(uint) (*dto.PhotoDto, error)
	PhotoList() (*[]dto.PhotoDto, error)
	CreatePhoto(photo []*dto.PhotoDto) []uint
	UpdatePhoto(photo *dto.PhotoDto) (*dto.PhotoDto, error)
	DeletePhoto(uint) (*dto.PhotoDto, error)
}

type photoService struct {
	ctx    application.AppContext
	db     *gorm.DB
	logger logger.AppLogger
}

func NewPhotoService(ctx application.AppContext, db *gorm.DB) PhotoService {
	return &photoService{ctx: ctx, db: db, logger: ctx.GetLogger()}
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

func (serv *photoService) saveUploadPhoto(photoDto *dto.PhotoDto) error {
	mf := photoDto.MultipartFile
	file, err := mf.Open()
	defer file.Close()
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		return err
	}

	sum := md5.Sum(buf.Bytes())
	hexSum := hex.EncodeToString(sum[:])

	var datas []model.Photo

	if result := serv.db.Where(&model.Photo{Sum: hexSum}).Find(&datas); result.Error != nil {
		if result.Error != nil {
			return result.Error
		}
	}
	if len(datas) > 0 {
		return errors.New("file exists")
	}

	photoDto.Sum = hexSum

	cfg, format, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	photoDto.Type = format
	photoDto.Resolution = fmt.Sprintf("%dx%d", cfg.Width, cfg.Height)

	fHandler := serv.ctx.GetFileHandler("")
	uri, err := fHandler.Save("", bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	photoDto.Size = mf.Size
	photoDto.Desc = fmt.Sprintf("%s\n%s", photoDto.Desc, mf.Filename)
	photoDto.Uri = uri

	photo := photoDto.ToModel()

	if result := serv.db.Create(photo); result.Error != nil {
		return result.Error
	}
	return nil
}

func (serv *photoService) CreatePhoto(photoDtos []*dto.PhotoDto) []uint {
	var failedUploadID []uint
	for _, dto := range photoDtos {
		if dto.MultipartFile != nil {
			err := serv.saveUploadPhoto(dto)
			if err != nil {
				serv.logger.Warn("file %s upload failed, %s", dto.MultipartFile.Filename, err.Error())
				failedUploadID = append(failedUploadID, dto.UploadID)
			}
		}
	}
	return failedUploadID
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
