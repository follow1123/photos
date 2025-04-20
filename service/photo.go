package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/filehandler"
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
	logger.AppLogger
	ctx    application.AppContext
	db     *gorm.DB
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

func (serv *photoService) saveUploadPhoto(photoDto *dto.PhotoDto) error {
	mf := photoDto.MultipartFile

	if !strings.Contains(photoDto.Desc, mf.Filename) {
		photoDto.Desc = fmt.Sprintf("%s\n%s", photoDto.Desc, mf.Filename)
	}
	file, err := mf.Open()
	if err != nil {
		serv.Error("open multipart file error: %s", err.Error())
		return err
	}
	defer file.Close()

	rootPath := filepath.Join(serv.ctx.GetConfig().GetFilesPath(), "test")
	imgProcessor, err := filehandler.NewImageProcessor(file)
	if err != nil {
		serv.Error("init image processor error: %s", err.Error())
		return err
	}

	// 判断数据库内是否存在相同的图片
	photoDto.Sum = imgProcessor.GetHexSum()
	result := serv.db.Select("id").Where(&model.Photo{Sum: photoDto.Sum}).Take(&model.Photo{});
	if result.Error == nil {
		return errors.New("file exists")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// 获取图片其他信息
	metaInfo, err := imgProcessor.GetMetaInfo()
	if err != nil {
		serv.Error("get image meta info error: %s", err.Error())
		return err
	}
	photoDto.Type = metaInfo.ImgFormat
	photoDto.Resolution = metaInfo.Resolution
	photoDto.Size = metaInfo.Size


	cacheFile, err := os.Create(filepath.Join(rootPath, "img_cache"))
	if err != nil {
		serv.Error("create cache file error: %s", err.Error())
		return err
	}
	defer cacheFile.Close()


	rawFile, err := os.Create(filepath.Join(rootPath, "img"))
	if err != nil {
		serv.Error("create raw file error: %s", err.Error())
		return err
	}
	defer rawFile.Close()

	_, err = io.Copy(rawFile, bytes.NewReader(imgProcessor.Data))
	if err != nil {
		serv.Error("save raw file error: %s", err.Error())
		return err
	}

	compressedData, _ := imgProcessor.GetCompressedData()
	_, err = io.Copy(cacheFile, bytes.NewReader(compressedData))
	if err != nil {
		serv.Error("save cache file error: %s", err.Error())
		return err
	}

	photo := photoDto.ToModel()
	if result := serv.db.Create(photo); result.Error != nil {
		return result.Error
	}
	return nil
}
	// fHandler := serv.ctx.GetFileHandler("")
	// uri, err := fHandler.Save("", bytes.NewReader(buf.Bytes()))
	// if err != nil {
	// 	return err
	// }
	// photoDto.Uri = uri


func (serv *photoService) CreatePhoto(photoDtos []*dto.PhotoDto) []uint {
	var failedUploadID []uint
	for _, dto := range photoDtos {
		if dto.MultipartFile != nil {
			err := serv.saveUploadPhoto(dto)
			if err != nil {
				serv.Warn("file %s upload failed, %s", dto.MultipartFile.Filename, err.Error())
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
