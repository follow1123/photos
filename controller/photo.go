package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/follow1123/photos/application"
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

type PhotoController struct {
	logger.AppLogger
	ctx  *application.AppContext
	serv service.PhotoService
}

func NewPhotoController(ctx *application.AppContext, service service.PhotoService) PhotoController {
	return PhotoController{ctx: ctx, serv: service, AppLogger: *ctx.GetLogger()}
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

func (pc *PhotoController) PhotoList(c *gin.Context) {
	data, err := pc.serv.PhotoList()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
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
				if uri := strings.TrimSpace(param.Uri); uri == "" {
					c.Error(application.NewAppError(http.StatusBadRequest, "不上传文件时，uri 必须填写"))
					return
				}
				continue
			}
			c.Error(application.NewAppError(http.StatusBadRequest, "上传了错误的文件"))
			return
		}
		param.FileHeader = fileHeader
	}
	if err := pc.serv.CreatePhoto(params); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
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
		extraHeaders["Content-Disposition"] = fmt.Sprintf(`attachment; filename="%s.%s"`, imgInfo.Name, imgInfo.Format)
	}
	c.DataFromReader(
		http.StatusOK,
		imgInfo.Size,
		fmt.Sprintf("image/%s", imgInfo.Format),
		rc,
		extraHeaders,
	)
}
