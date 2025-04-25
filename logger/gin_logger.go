package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const GIN_PREFIX = "[GIN]"

type GinLogger struct {
	Logger *zap.SugaredLogger
}

func NewGinLogger(baseLogger *zap.SugaredLogger) *GinLogger {
	return &GinLogger{
		Logger: baseLogger.WithOptions(zap.WithCaller(false)).Named(GIN_PREFIX),
	}
}

func (gl *GinLogger) Handler(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	// Process request
	c.Next()

	// Stop timer
	latency := time.Now().Sub(start)
	if latency > time.Minute {
		latency = latency.Truncate(time.Second)
	}

	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

	if raw != "" {
		path = path + "?" + raw
	}

	if errorMessage != "" {
		errorMessage = "\n" + errorMessage
	}
	gl.Logger.Debugf("%3d | %13v | %15s | %-7s  %#v%s",
		statusCode,
		latency,
		clientIP,
		method,
		path,
		errorMessage,
	)

}
