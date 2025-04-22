package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/common"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/model/dto"
	"github.com/follow1123/photos/service"
	"github.com/gin-gonic/gin"
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

type PhotoController interface {
	GetPhotoById(*gin.Context)
	PhotoList(*gin.Context)
	CreatePhoto(*gin.Context)
	UpdatePhoto(*gin.Context)
	DeletePhoto(*gin.Context)
	PreviewOriginalPhoto(*gin.Context)
}

type photoController struct {
	logger.AppLogger
	ctx     application.AppContext
	service service.PhotoService
}

func NewPhotoController(ctx application.AppContext, service service.PhotoService) PhotoController {
	return &photoController{ctx: ctx, service: service, AppLogger: *ctx.GetLogger()}
}

func (ctl *photoController) GetPhotoById(c *gin.Context) {
	data := &dto.PhotoDto{}
	var err error
	if err = c.BindUri(data); err != nil {
		return
	}
	data, err = ctl.service.GetPhotoById(data.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if data == nil {
		c.AbortWithError(http.StatusNotFound, common.ErrDataNotFound).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, data)
}

func (ctl *photoController) PhotoList(c *gin.Context) {
	data, err := ctl.service.PhotoList()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if data == nil {
		c.AbortWithError(http.StatusNotFound, common.ErrDataNotFound).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, data)
}

func (ctl *photoController) CreatePhoto(c *gin.Context) {
	metaData := c.PostForm("metaData")

	ctl.Info("meta data: %s", metaData)
	var metaDatas []*dto.PhotoDto
	if err := json.Unmarshal([]byte(metaData), &metaDatas); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid meta data")).SetType(gin.ErrorTypePublic)
		return
	}

	uploadFiles := make(map[uint]*multipart.FileHeader, len(metaDatas))

	for _, data := range metaDatas {
		fileHeader, err := c.FormFile(fmt.Sprintf("file_%d", data.UploadID))
		missingFile := errors.Is(err, http.ErrMissingFile)
		if missingFile {
			uri := strings.TrimSpace(data.Uri)
			if uri == "" {
				c.AbortWithError(http.StatusBadRequest, errors.New("when no file is uploaded, uri is required")).SetType(gin.ErrorTypePublic)
				return
			}
			continue
		}
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("upload a bad file")).SetType(gin.ErrorTypePublic)
			return
		}
		uploadFiles[data.UploadID] = fileHeader
	}
	c.JSON(http.StatusOK, ctl.service.CreatePhoto(metaDatas, uploadFiles))
}

func (ctl *photoController) UpdatePhoto(c *gin.Context) {
	data := &dto.PhotoDto{Operate: dto.PHOTO_OPE_ID_REQUIRED}
	var err error
	if err = c.BindJSON(&data); err != nil {
		return
	}
	if data, err = ctl.service.UpdatePhoto(data); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if data == nil {
		c.AbortWithError(http.StatusNotFound, common.ErrDataNotFound).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (ctl *photoController) DeletePhoto(c *gin.Context) {
	data := &dto.PhotoDto{Operate: dto.PHOTO_OPE_ID_REQUIRED}
	var err error
	if err = c.BindUri(data); err != nil {
		return
	}
	if data, err = ctl.service.DeletePhoto(data.ID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if data == nil {
		c.AbortWithError(http.StatusNotFound, common.ErrDataNotFound).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (ctl *photoController) PreviewOriginalPhoto(c *gin.Context) {
	var photoDto dto.PhotoDto
	var err error
	if err = c.BindUri(&photoDto); err != nil {
		return
	}

	urlPath := c.Request.URL.Path
	isDownload := strings.HasSuffix(urlPath, "download")
	isCompressed := strings.HasSuffix(urlPath, "compressed")

	rc, err := ctl.service.GetPhotoFile(&photoDto, !isCompressed)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if rc == nil {
		c.AbortWithError(http.StatusNotFound, common.ErrDataNotFound).SetType(gin.ErrorTypePublic)
		return
	}
	defer rc.Close()
    extraHeaders := make(map[string]string, 1)

	var fileName string = "qwrewqrwqe"

	if isDownload {
		extraHeaders["Content-Disposition"] = fmt.Sprintf(`attachment; filename="%s.%s"`, fileName, photoDto.Format)
	}
	c.DataFromReader(
		http.StatusOK,
		photoDto.Size,
		fmt.Sprintf("image/%s", photoDto.Format),
		rc,
		extraHeaders,
	)
}
