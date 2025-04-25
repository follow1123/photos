package webserver

import (
	"errors"

	"github.com/follow1123/photos/application"
	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type Router interface {
	SetHandleMapping(engine *gin.Engine)
}

type GinWebServer struct {
	engine  *gin.Engine
	logger  *logger.GinLogger
	conf    *config.Config
	routers []Router
}

func NewGinWebServer(config *config.Config, ginLogger *logger.GinLogger) *GinWebServer {
	return &GinWebServer{
		engine:  gin.New(),
		conf:    config,
		logger:  ginLogger,
		routers: make([]Router, 0, 0),
	}
}

func (gws *GinWebServer) UseLoggerMiddleware() {
	gws.engine.Use(gws.logger.Handler)
}

func (gws *GinWebServer) UseRecoveryMiddleware() {
	gws.engine.Use(gin.Recovery())
}

func (gws *GinWebServer) UseErrorHandlerMiddleware() {
	gws.engine.Use(gws.globalErrorHandler)
}

func (gws *GinWebServer) InitMiddleware() {
	gws.UseLoggerMiddleware()
	gws.UseRecoveryMiddleware()
	gws.UseErrorHandlerMiddleware()
}

func (gws *GinWebServer) Start() {
	gws.engine.Run(gws.conf.GetAddr())
}

func (gws *GinWebServer) globalErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) == 0 {
		return
	}

	err := c.Errors[0]

	for _, e := range c.Errors {
		// 优先处理内部错误
		if e.IsType(gin.ErrorTypePrivate) {
			err = e
			break
		}
	}

	if err.IsType(gin.ErrorTypePrivate) {
		var appError *application.AppError
		if !errors.As(err, &appError) {
			appError = application.ErrInternalServerError
		}
		c.JSON(appError.Code, appError)
	} else {
		r := render.JSON{Data: application.AppError{
			Message: err.Error(),
		}}
		if err := r.Render(c.Writer); err != nil {
			gws.logger.Logger.Error("render message error")
		}
	}
}

func (gws *GinWebServer) SetRouters(routers ...Router) {
	gws.routers = append(gws.routers, routers...)
}

func (gws *GinWebServer) InitRouter() {
	l := gws.logger.Logger
	l.Infof("init %d gin router", len(gws.routers))
	for _, router := range gws.routers {
		router.SetHandleMapping(gws.engine)
	}
}

func (gws *GinWebServer) GetEngine() *gin.Engine {
	return gws.engine
}
