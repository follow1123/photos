package controller

import (
	"net/http"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/common"
	"github.com/follow1123/photos/logger"
	"github.com/follow1123/photos/model/dto"
	"github.com/follow1123/photos/service"
	"github.com/gin-gonic/gin"
)

const (
	PHOTO_API_GETBYID string = "/photo/:id"
	PHOTO_API_LIST           = "/photo"
	PHOTO_API_CREATE         = PHOTO_API_LIST
	PHOTO_API_UPDATE         = PHOTO_API_LIST
	PHOTO_API_DELETE         = PHOTO_API_GETBYID
)

type PhotoController interface {
	GetPhotoById(*gin.Context)
	PhotoList(*gin.Context)
	CreatePhoto(*gin.Context)
	UpdatePhoto(*gin.Context)
	DeletePhoto(*gin.Context)
}

type photoController struct {
	logger  logger.AppLogger
	ctx     application.AppContext
	service service.PhotoService
}

func NewPhotoController(ctx application.AppContext, service service.PhotoService) PhotoController {
	return &photoController{logger: ctx.GetLogger(), ctx: ctx, service: service}
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
	data := &dto.PhotoDto{Operate: dto.PHOTO_OPE_CREATE}
	var err error

	if err = c.BindJSON(&data); err != nil {
		return
	}
	if data, err = ctl.service.CreatePhoto(data); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if data == nil {
		c.AbortWithError(http.StatusNotFound, common.ErrDataNotFound).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, data)
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
