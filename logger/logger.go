package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitBaseLogger() *zap.SugaredLogger {
	zapConf := zap.NewDevelopmentConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	logger, err := zapConf.Build()
	logger.WithOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "init logger error")
		os.Exit(1)
	}
	return logger.Sugar()
}
