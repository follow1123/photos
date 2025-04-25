package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewBaseLogger() (*zap.SugaredLogger, error) {
	zapConf := zap.NewDevelopmentConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	logger, err := zapConf.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
