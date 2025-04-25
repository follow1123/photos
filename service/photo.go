package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/database"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/model"
	"github.com/follow1123/photos/model/dto"
	"gorm.io/gorm"
)

type PhotoService interface {
	GetPhotoById(uint) (*dto.PhotoDto, error)
	PhotoPage(dto.PageParam[dto.PhotoPageParam]) (*dto.PageResult[dto.PhotoDto], error)
	CreatePhoto([]dto.CreatePhotoParam) []dto.CreatePhotoFailedResult
	UpdatePhoto(dto.PhotoParam) (*dto.PhotoDto, error)
	DeletePhoto(uint) error
	GetPhotoFile(uint, bool) (io.ReadCloser, *imagemanager.ImageInfo, error)
}

type photoService struct {
	logger.AppLogger
	ctx *application.AppContext
	db  *database.SqliteDB
}

func NewPhotoService(ctx *application.AppContext, db *database.SqliteDB) PhotoService {
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

func (ps *photoService) PhotoPage(pageParam dto.PageParam[dto.PhotoPageParam]) (*dto.PageResult[dto.PhotoDto], error) {
	var photoDtoList []dto.PhotoDto
	result := ps.db.Table("photos").Where(pageParam.Params.ToModel()).Offset(pageParam.PageNum - 1).Limit(pageParam.PageSize).Find(&photoDtoList)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(photoDtoList) == 0 {
		return nil, application.ErrDataNotFound
	}

	return &dto.PageResult[dto.PhotoDto]{
		List:     photoDtoList,
		PageNum:  pageParam.PageNum,
		PageSize: pageParam.PageSize,
		Total:    0,
	}, nil

}

// func (ps *photoService) saveUploadPhoto(param *dto.CreatePhotoParam) error {
// 	uploadMgr := ps.ctx.GetImageManager().NewUploadManager(
// 		imagemanager.NewMultipartSource(param.FileHeader),
// 	)
// 	photo := param.ToModel()
// 	// 判断数据库内是否存在相同的图片
//
// 	sum, err := uploadMgr.GetHexSum()
// 	if err != nil {
// 		return err
// 	}
//
// 	result := ps.db.Select("id").Where(&model.Photo{Sum: sum}).Take(&model.Photo{})
// 	if result.Error == nil {
// 		return errors.New("file exists")
// 	}
// 	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return result.Error
// 	}
//
// 	photo.Sum = sum
//
// 	// 获取图片其他信息
// 	imgInfo, err := uploadMgr.GetImageInfo()
// 	if err != nil {
// 		ps.Error("get image meta info error: %s", err.Error())
// 		return err
// 	}
//
// 	imageName := uploadMgr.GetImageName()
// 	if !strings.Contains(photo.Desc, imageName) {
// 		photo.Desc = fmt.Sprintf("%s\n%s", photo.Desc, imageName)
// 	}
// 	photo.Size = imgInfo.Size
// 	photo.Format = imgInfo.Format
// 	photo.Width = imgInfo.Width
// 	photo.Height = imgInfo.Height
//
// 	uri, err := uploadMgr.Save()
// 	if err != nil {
// 		return err
// 	}
// 	photo.Uri = uri
//
// 	if result := ps.db.Create(photo); result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

func (ps *photoService) CreatePhoto(params []dto.CreatePhotoParam) []dto.CreatePhotoFailedResult {
	var (
		preparedPhotos = make([]model.Photo, 0, len(params))
		failedResults  = make([]dto.CreatePhotoFailedResult, 0, len(params))
	)
	var (
		numWorkers  = 8
		numJobs     = len(params)
		numModels   = len(params)
		numFailures = len(params)
	)
	var (
		wg     sync.WaitGroup
		sumMap sync.Map
	)

	jobs := make(chan int, numJobs)
	models := make(chan *model.Photo, numModels)
	failures := make(chan *dto.CreatePhotoFailedResult, numFailures)

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				param := params[job]
				uploadMgr := ps.ctx.GetImageManager().NewUploadManager(param.ImageSource)
				var photo = model.Photo{Desc: param.Desc, PhotoDate: param.PhotoDate}
				imageName := uploadMgr.GetImageName()
				if !strings.Contains(photo.Desc, imageName) {
					photo.Desc = ConcatDesc(photo.Desc, imageName)
				}

				sum, err := uploadMgr.GetHexSum()
				if err != nil {
					ps.Error("get hex sum error: %v", err)
					failures <- &dto.CreatePhotoFailedResult{
						UploadID: param.UploadID,
						Message:  err.Error(),
					}
					continue
				}

				// 判断是否和其他正在上传的文件重复
				_, loaded := sumMap.LoadOrStore(sum, true)
				if loaded {
					failures <- &dto.CreatePhotoFailedResult{
						UploadID: param.UploadID,
						Message:  BuildUploadDupMsg(imageName, photo.Desc),
					}
					continue
				}

				// 判断数据库内是否存在相同的图片
				result := ps.db.Select("id").Where(&model.Photo{Sum: sum}).Take(&model.Photo{})
				if result.Error == nil {
					msg := "文件重复"
					ps.Error(msg)
					failures <- &dto.CreatePhotoFailedResult{
						UploadID: param.UploadID,
						Message:  msg,
					}
					continue
				}
				if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
					ps.Error("select same sum photo error: %v", err)
					failures <- &dto.CreatePhotoFailedResult{
						UploadID: param.UploadID,
						Message:  result.Error.Error(),
					}
					continue
				}

				photo.Sum = sum

				// 获取图片其他信息
				imgInfo, err := uploadMgr.GetImageInfo()
				if err != nil {
					ps.Error("get image info error: %v", err)
					failures <- &dto.CreatePhotoFailedResult{
						UploadID: param.UploadID,
						Message:  err.Error(),
					}
					continue
				}

				photo.Size = imgInfo.Size
				photo.Format = imgInfo.Format
				photo.Width = imgInfo.Width
				photo.Height = imgInfo.Height

				// 保存图片
				uri, err := uploadMgr.Save()
				if err != nil {
					ps.Error("save image error: %v", err)
					failures <- &dto.CreatePhotoFailedResult{
						UploadID: param.UploadID,
						Message:  err.Error(),
					}
					continue
				}
				photo.Uri = uri

				// 加入待保存列表
				models <- &photo
			}
		}()
	}

	for j := range numJobs {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(models)
	close(failures)
	for failedResult := range failures {
		failedResults = append(failedResults, *failedResult)
	}

	for photo := range models {
		preparedPhotos = append(preparedPhotos, *photo)
	}

	if result := ps.db.Create(&preparedPhotos); result.Error != nil {
		panic(fmt.Sprintf("batch save photo error: %v", result.Error))
	}

	return failedResults
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

	imgInfo := &imagemanager.ImageInfo{Size: photo.Size, Format: photo.Format}
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

func ConcatDesc(desc string, name string) string {
	return fmt.Sprintf("%s\n%s", desc, name)
}

func BuildUploadDupMsg(dup string, last string) string {
	return fmt.Sprintf("上传的文件内 [ %s ] 和 [ %s（已保存）] 重复", dup, last)
}
