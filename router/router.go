package router

import (
	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/controller"
	"github.com/follow1123/photos/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const LOG_PREFIX = "[ROUTER]"

func Init(r *gin.Engine, appCtx *application.AppContext, baseLogger *zap.SugaredLogger, db *gorm.DB) {
	baseLogger.Debugf("%s init router", LOG_PREFIX)

	photoCtl := controller.NewPhotoController(appCtx, service.NewPhotoService(appCtx, db))
	r.GET(controller.PHOTO_API_GETBYID, photoCtl.GetPhotoById)
	r.GET(controller.PHOTO_API_LIST, photoCtl.PhotoList)
	r.POST(controller.PHOTO_API_CREATE, photoCtl.CreatePhoto)
	r.PUT(controller.PHOTO_API_UPDATE, photoCtl.UpdatePhoto)
	r.DELETE(controller.PHOTO_API_DELETE, photoCtl.DeletePhoto)
	r.GET(controller.PHOTO_API_PREVIEW_ORIGINAL, photoCtl.PreviewOriginalPhoto)
	r.GET(controller.PHOTO_API_PREVIEW_COMPRESSED, photoCtl.PreviewOriginalPhoto)
	r.GET(controller.PHOTO_API_DOWNLOAD, photoCtl.PreviewOriginalPhoto)
}
