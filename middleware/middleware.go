package middleware

import (
	"errors"

	"github.com/follow1123/photos/application"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"go.uber.org/zap"
)

const ERROR_HANDLER_PREFIX = "[ERROR_HANDLER]"

func GlobalErrorHandler(baseLogger *zap.SugaredLogger) gin.HandlerFunc {
	errLogger := baseLogger.Named(ERROR_HANDLER_PREFIX)
	return func(c *gin.Context) {
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
				errLogger.Error("render message error")
			}
		}
	}
}
