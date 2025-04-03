package logger

import (
	"go.uber.org/zap"
)

const APP_PREFIX = "[APP]"

type AppLogger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	Fatal(string, ...any)
}

type appLogger struct {
	logger *zap.SugaredLogger
}

func NewAppLogger(baseLogger *zap.SugaredLogger) AppLogger {
	return &appLogger{logger: baseLogger.WithOptions(zap.AddCallerSkip(1)).Named(APP_PREFIX)}
}

func (al *appLogger) Debug(template string, args ...any) {
	al.logger.Debugf(template, args...)
}

func (al *appLogger) Info(template string, args ...any) {
	al.logger.Infof(template, args...)
}

func (al *appLogger) Warn(template string, args ...any) {
	al.logger.Warnf(template, args...)
}

func (al *appLogger) Error(template string, args ...any) {
	al.logger.Errorf(template, args...)
}

func (al *appLogger) Fatal(template string, args ...any) {
	al.logger.Fatalf(template, args...)
}
