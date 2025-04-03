package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
	// gormUtils "gorm.io/gorm/utils"
)

const (
	GORM_PREFIX   = "[GORM] "
	slowThreshold = 200
)

type GormLogger struct {
	logger *zap.SugaredLogger
}

func NewGormLogger(baseLogger *zap.SugaredLogger) gormLogger.Interface {
	return &GormLogger{logger: baseLogger.WithOptions(zap.AddCallerSkip(3)).Named(GORM_PREFIX)}
}
func (gl *GormLogger) LogMode(_ gormLogger.LogLevel) gormLogger.Interface {
	return gl
}

func (gl *GormLogger) Info(_ context.Context, msg string, data ...any) {
	gl.logger.Infof(msg, data...)
}

func (gl *GormLogger) Warn(_ context.Context, msg string, data ...any) {
	gl.logger.Warnf(msg, data...)
}

func (gl *GormLogger) Error(_ context.Context, msg string, data ...any) {
	gl.logger.Errorf(msg, data...)
}

func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// if gl.logger.Level() != zapcore.DebugLevel {
	// 	return
	// }
	elapsed := time.Since(begin)
	switch {
	case err != nil:
		sql, rows := fc()
		if rows == -1 {
			gl.logger.Errorf("%s\n[%.3fms] [rows:%v] %s\n", err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			gl.logger.Errorf("%s\n[%.3fms] [rows:%v] %s\n", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > slowThreshold*time.Millisecond && slowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", slowThreshold*time.Millisecond)
		if rows == -1 {
			gl.logger.Warnf("%s\n[%.3fms] [rows:%v] %s\n", slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			gl.logger.Warnf("%s\n[%.3fms] [rows:%v] %s\n", slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			gl.logger.Debugf("\n[%.3fms] [rows:%v] %s\n", float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			gl.logger.Debugf("\n[%.3fms] [rows:%v] %s\n", float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
