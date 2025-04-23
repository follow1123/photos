package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/model"
	"github.com/follow1123/photos/model/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PhotoService interface {
	GetPhotoById(uint) (*dto.PhotoDto, error)
	PhotoList() (*[]dto.PhotoDto, error)
	CreatePhoto([]dto.CreatePhotoParam) error
	UpdatePhoto(dto.PhotoParam) (*dto.PhotoDto, error)
	DeletePhoto(uint) error
	GetPhotoFile(uint, bool) (io.ReadCloser, *imagemanager.ImageInfo, error)
}

type photoService struct {
	logger.AppLogger
	ctx *application.AppContext
	db  *gorm.DB
}

func NewPhotoService(ctx *application.AppContext, db *gorm.DB) PhotoService {
	return &photoService{ctx: ctx, db: db, AppLogger: *ctx.GetLogger()}
}

func (ps *photoService) GetPhotoById(id uint) (*dto.PhotoDto, error) {
	var photo model.Photo
	if result := ps.db.First(&photo, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, application.ErrDataNotFound
		}
		return nil, result.Error
	}
	photoDto := &dto.PhotoDto{}
	photoDto.Update(&photo)
	return photoDto, nil
}

func (ps *photoService) PhotoList() (*[]dto.PhotoDto, error) {
	var photoList []model.Photo

	if result := ps.db.Find(&photoList); result.Error != nil {
		return nil, result.Error
	}

	if len(photoList) == 0 {
		return nil, application.ErrDataNotFound
	}

	dataList := make([]dto.PhotoDto, len(photoList))

	for i, v := range photoList {
		dataList[i] = dto.PhotoDto{}
		dataList[i].Update(&v)
	}

	return &dataList, nil
}

func (ps *photoService) saveUploadPhoto(param *dto.CreatePhotoParam) error {
	imgManager, err := ps.ctx.GetImageManager().NewUploadManager(
		&imagemanager.MultipartSource{FileHeader: param.FileHeader},
	)
	if err != nil {
		return err
	}
	photo := param.ToModel()

	// 判断数据库内是否存在相同的图片
	photo.Sum = imgManager.GetHexSum()
	result := ps.db.Select("id").Where(&model.Photo{Sum: photo.Sum}).Take(&model.Photo{})
	if result.Error == nil {
		return errors.New("file exists")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// 获取图片其他信息
	imgInfo, err := imgManager.GetImageInfo()
	if err != nil {
		ps.Error("get image meta info error: %s", err.Error())
		return err
	}
	if !strings.Contains(photo.Desc, imgInfo.Name) {
		photo.Desc = fmt.Sprintf("%s\n%s", photo.Desc, imgInfo.Name)
	}
	photo.Size = imgInfo.Size
	photo.Format = imgInfo.Format
	photo.Width = imgInfo.Width
	photo.Height = imgInfo.Height

	uri, err := imgManager.Save()
	if err != nil {
		return err
	}
	photo.Uri = uri

	if result := ps.db.Create(photo); result.Error != nil {
		return result.Error
	}
	return nil
}

func (ps *photoService) CreatePhoto(params []dto.CreatePhotoParam) error {
	var failedResults []dto.CreatePhotoFailedResult
	for _, p := range params {
		if p.FileHeader != nil {
			if err := ps.saveUploadPhoto(&p); err != nil {
				failedResults = append(failedResults, dto.CreatePhotoFailedResult{
					UploadID: p.UploadID,
					Message:  err.Error(),
				})
			}
		}
	}
	failureCount := len(failedResults)
	if failureCount > 0 {
		successCount := len(params) - failureCount
		return application.NewAppError(
			http.StatusMultiStatus,
			"%d 个文件上传成功，%d 个文件上传失败",
			successCount,
			failureCount,
		)
	}

	return nil
}

func (ps *photoService) UpdatePhoto(param dto.PhotoParam) (*dto.PhotoDto, error) {
	var photo model.Photo
	if result := ps.db.First(&photo, param.ID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, application.ErrDataNotFound
		}
		return nil, result.Error
	}
	if result := ps.db.Model(&photo).Updates(param); result.Error != nil {
		return nil, result.Error
	}
	photoDto := &dto.PhotoDto{}
	photoDto.Update(&photo)
	return photoDto, nil
}

func (ps *photoService) DeletePhoto(id uint) error {
	photo := &model.Photo{}
	if result := ps.db.First(photo, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return application.ErrDataNotFound
		}
		return result.Error
	}
	if result := ps.db.Delete(photo, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (ps *photoService) GetPhotoFile(id uint, original bool) (io.ReadCloser, *imagemanager.ImageInfo, error) {
	var photo model.Photo
	if result := ps.db.First(&photo, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil, application.ErrDataNotFound
		}
		return nil, nil, result.Error
	}

	downloadManager := ps.ctx.GetImageManager().NewDownloadManager(photo.Uri)

	t := time.Now()
	timestamp := t.Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s", timestamp, strings.ReplaceAll(uuid.New().String(), "-", ""))
	imgInfo := &imagemanager.ImageInfo{Name: fileName, Size: photo.Size}
	if original {
		reader, err := downloadManager.OpenOriginal()
		if err != nil {
			return nil, nil, err
		}
		return reader, imgInfo, nil
	}

	imageData, err := downloadManager.GetCompressed()
	if err != nil {
		return nil, nil, err
	}
	imgInfo.Size = int64(len(imageData))
	return io.NopCloser(bytes.NewReader(imageData)), imgInfo, nil
}
