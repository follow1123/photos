package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/model"
	"github.com/follow1123/photos/model/dto"
	"gorm.io/gorm"
)

type PhotoService interface {
	GetPhotoById(uint) (*dto.PhotoDto, error)
	PhotoList() (*[]dto.PhotoDto, error)
	CreatePhoto(photo []*dto.PhotoDto, uploadFiles map[uint]*multipart.FileHeader) []dto.CreatePhotoFailedResult
	UpdatePhoto(photo *dto.PhotoDto) (*dto.PhotoDto, error)
	DeletePhoto(uint) (*dto.PhotoDto, error)
	GetPhotoFile(*dto.PhotoDto, bool) (io.ReadCloser, error)
}

type photoService struct {
	logger.AppLogger
	ctx application.AppContext
	db  *gorm.DB
}

func NewPhotoService(ctx application.AppContext, db *gorm.DB) PhotoService {
	return &photoService{ctx: ctx, db: db, AppLogger: *ctx.GetLogger()}
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

func (serv *photoService) saveUploadPhoto(photoDto *dto.PhotoDto, fileHeader *multipart.FileHeader) error {
	imgManager, err := serv.ctx.GetImageManager().NewUploadManager(
		&imagemanager.MultipartSource{FileHeader: fileHeader},
	)
	if err != nil {
		return err
	}

	// 判断数据库内是否存在相同的图片
	photoDto.Sum = imgManager.GetHexSum()
	result := serv.db.Select("id").Where(&model.Photo{Sum: photoDto.Sum}).Take(&model.Photo{})
	if result.Error == nil {
		return errors.New("file exists")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// 获取图片其他信息
	imgInfo, err := imgManager.GetImageInfo()
	if err != nil {
		serv.Error("get image meta info error: %s", err.Error())
		return err
	}
	if !strings.Contains(photoDto.Desc, imgInfo.Name) {
		photoDto.Desc = fmt.Sprintf("%s\n%s", photoDto.Desc, imgInfo.Name)
	}
	photoDto.Size = imgInfo.Size
	photoDto.Format = imgInfo.Format
	photoDto.Width = imgInfo.Width
	photoDto.Height = imgInfo.Height

	uri, err := imgManager.Save()
	if err != nil {
		return err
	}
	photoDto.Uri = uri

	photo := photoDto.ToModel()
	if result := serv.db.Create(photo); result.Error != nil {
		return result.Error
	}
	return nil
}

func (serv *photoService) CreatePhoto(photoDtos []*dto.PhotoDto, uploadFiles map[uint]*multipart.FileHeader) []dto.CreatePhotoFailedResult {
	var failedResults []dto.CreatePhotoFailedResult
	for _, p := range photoDtos {
		fileHeader, ok := uploadFiles[p.UploadID]
		if ok {
			if err := serv.saveUploadPhoto(p, fileHeader); err != nil {
				failedResults = append(failedResults, dto.CreatePhotoFailedResult{
					UploadID: p.UploadID,
					Message:  err.Error(),
				})
			}
		}
	}
	return failedResults
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

func (serv *photoService) GetPhotoFile(photoDto *dto.PhotoDto, original bool) (io.ReadCloser, error) {
	var photo model.Photo
	if result := serv.db.First(&photo, photoDto.ID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	photoDto.Update(&photo)

	downloadManager := serv.ctx.GetImageManager().NewDownloadManager(photo.Uri)
	if original {
		return downloadManager.OpenOriginal()
	}

	imageData, err := downloadManager.GetCompressed()
	if err != nil {
		return nil, err
	}
	photoDto.Size = int64(len(imageData))
	return io.NopCloser(bytes.NewReader(imageData)), nil
}

