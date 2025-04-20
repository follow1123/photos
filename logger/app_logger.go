package logger

import (
	"go.uber.org/zap"
)

const APP_PREFIX = "[APP]"

type AppLogger struct {
	Logger *zap.SugaredLogger
}

func NewAppLogger(baseLogger *zap.SugaredLogger) *AppLogger {
	return &AppLogger{baseLogger.WithOptions(zap.AddCallerSkip(1)).Named(APP_PREFIX)}
}

func (w *AppLogger) Debug(template string, args ...any) {
	w.Logger.Debugf(template, args...)
}

func (w *AppLogger) Info(template string, args ...any) {
	w.Logger.Infof(template, args...)
}

func (w *AppLogger) Warn(template string, args ...any) {
	w.Logger.Warnf(template, args...)
}

func (w *AppLogger) Error(template string, args ...any) {
	w.Logger.Errorf(template, args...)
}

func (w *AppLogger) Fatal(template string, args ...any) {
	w.Logger.Fatalf(template, args...)
}
