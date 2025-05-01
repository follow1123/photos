package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/model/dto"
	"github.com/follow1123/photos/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	PHOTO_API_GETBYID            string = "/photo/:id"
	PHOTO_API_LIST                      = "/photo"
	PHOTO_API_CREATE                    = PHOTO_API_LIST
	PHOTO_API_UPDATE                    = PHOTO_API_LIST
	PHOTO_API_DELETE                    = PHOTO_API_GETBYID
	PHOTO_API_PREVIEW_ORIGINAL          = PHOTO_API_GETBYID + "/preview/original"
	PHOTO_API_PREVIEW_COMPRESSED        = PHOTO_API_GETBYID + "/preview/compressed"
	PHOTO_API_DOWNLOAD                  = PHOTO_API_GETBYID + "/download"
)

type PhotoController struct {
	logger.AppLogger
	ctx  *application.AppContext
	serv service.PhotoService
}

func NewPhotoController(ctx *application.AppContext, service service.PhotoService) *PhotoController {
	return &PhotoController{ctx: ctx, serv: service, AppLogger: *ctx.GetLogger()}
}

func (pc *PhotoController) GetPhotoById(c *gin.Context) {
	param := &dto.PhotoParam{}
	if err := c.BindUri(param); err != nil {
		return
	}
	photoDto, err := pc.serv.GetPhotoById(param.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, photoDto)
}

func (pc *PhotoController) PhotoPage(c *gin.Context) {
	pageParam := dto.PageParam[dto.PhotoPageParam]{}
	photoParam := dto.PhotoPageParam{}

	if err := c.BindQuery(&pageParam); err != nil {
		return
	}
	if err := c.BindQuery(&photoParam); err != nil {
		return
	}

	result, err := pc.serv.PhotoPage(pageParam)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (pc *PhotoController) CreatePhoto(c *gin.Context) {
	metaData := c.PostForm("metaData")
	pc.Debug("meta data: %s", metaData)

	var params []dto.CreatePhotoParam
	if err := json.Unmarshal([]byte(metaData), &params); err != nil {
		c.Error(application.NewAppError(http.StatusBadRequest, "解析 metaData 的 json 数据错误: %v", err))
		return
	}

	for _, param := range params {
		fileHeader, err := c.FormFile(fmt.Sprintf("file_%d", param.UploadID))
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				param.Uri = strings.TrimSpace(param.Uri)
				if param.Uri == "" {
					c.Error(application.NewAppError(http.StatusBadRequest, "不上传文件时，uri 必须填写"))
					return
				}
				param.ImageSource = imagemanager.NewRemoteUriSource(param.Uri)
				continue
			}
			c.Error(application.NewAppError(http.StatusBadRequest, "上传了错误的文件"))
			return
		}
		param.ImageSource = imagemanager.NewMultipartSource(fileHeader)
	}
	failedResults := pc.serv.CreatePhoto(params)
	failureCount := len(failedResults)
	if failureCount == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	successCount := len(params) - failureCount
	appError := application.NewAppError(
		http.StatusMultiStatus,
		"%d 个文件上传成功，%d 个文件上传失败",
		successCount,
		failureCount,
	)
	appError.Details = failedResults
	c.Error(appError)
}

func (pc *PhotoController) UpdatePhoto(c *gin.Context) {
	var param dto.PhotoParam
	if err := c.BindJSON(&param); err != nil {
		return
	}
	photoDto, err := pc.serv.UpdatePhoto(param)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, photoDto)
}

func (pc *PhotoController) DeletePhoto(c *gin.Context) {
	param := &dto.PhotoParam{}
	if err := c.BindUri(param); err != nil {
		return
	}
	err := pc.serv.DeletePhoto(param.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (pc *PhotoController) PreviewOriginalPhoto(c *gin.Context) {
	param := &dto.PhotoParam{}
	if err := c.BindUri(param); err != nil {
		return
	}

	urlPath := c.Request.URL.Path
	isDownload := strings.HasSuffix(urlPath, "download")
	isCompressed := strings.HasSuffix(urlPath, "compressed")

	rc, imgInfo, err := pc.serv.GetPhotoFile(param.ID, !isCompressed)
	if err != nil {
		c.Error(err)
		return
	}
	defer rc.Close()
	extraHeaders := make(map[string]string, 1)

	if isDownload {
		t := time.Now()
		timestamp := t.Format("20060102150405")
		fileName := fmt.Sprintf("%s_%s", timestamp, strings.ReplaceAll(uuid.New().String(), "-", ""))
		extraHeaders["Content-Disposition"] = fmt.Sprintf(`attachment; filename="%s.%s"`, fileName, imgInfo.Format)
	}
	c.DataFromReader(
		http.StatusOK,
		imgInfo.Size,
		fmt.Sprintf("image/%s", imgInfo.Format),
		rc,
		extraHeaders,
	)
}

func (pc *PhotoController) SetHandleMapping(engine *gin.Engine) {
	engine.GET(PHOTO_API_GETBYID, pc.GetPhotoById)
	engine.GET(PHOTO_API_LIST, pc.PhotoPage)
	engine.POST(PHOTO_API_CREATE, pc.CreatePhoto)
	engine.PUT(PHOTO_API_UPDATE, pc.UpdatePhoto)
	engine.DELETE(PHOTO_API_DELETE, pc.DeletePhoto)
	engine.GET(PHOTO_API_PREVIEW_ORIGINAL, pc.PreviewOriginalPhoto)
	engine.GET(PHOTO_API_PREVIEW_COMPRESSED, pc.PreviewOriginalPhoto)
	engine.GET(PHOTO_API_DOWNLOAD, pc.PreviewOriginalPhoto)
}
