package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"go.uber.org/zap"
)

const ERROR_HANDLER_PREFIX = "[ERROR_HANDLER]"

type ErrorMsg struct {
	Message string `json:"message"`
}

func renderErrorMsg(c *gin.Context, message ErrorMsg, logger *zap.SugaredLogger) {
	r := render.JSON{Data: message}
	if err := r.Render(c.Writer); err != nil {
		logger.Error("render message error")
	}
}

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
			renderErrorMsg(c, ErrorMsg{Message: "Internal Server Error"}, errLogger)
		} else {
			renderErrorMsg(c, ErrorMsg{Message: err.Error()}, errLogger)
		}
	}
}
