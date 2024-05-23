package logger

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	olog "gorm.io/gorm/logger"
)

type ContextFn func(ctx context.Context) []zapcore.Field

type GormLogger struct {
	LogLevel                  olog.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	Context                   ContextFn
}

func DbLogger() GormLogger {
	var level olog.LogLevel
	switch handler.logger.Level() {
	case zap.WarnLevel:
		level = olog.Warn
	case zap.ErrorLevel:
		level = olog.Error
	case zap.DebugLevel:
		level = olog.Info
	default:
		level = olog.Silent
	}
	return GormLogger{
		LogLevel:                  level,
		SlowThreshold:             100 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
		Context:                   nil,
	}
}

func (l GormLogger) LogMode(level olog.LogLevel) olog.Interface {
	return GormLogger{
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
		Context:                   l.Context,
	}
}

func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < olog.Info {
		return
	}
	l.logger(ctx).Sugar().Debugf(str, args...)
}

func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < olog.Warn {
		return
	}
	l.logger(ctx).Sugar().Warnf(str, args...)
}

func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < olog.Error {
		return
	}
	l.logger(ctx).Sugar().Errorf(str, args...)
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	logger := l.logger(ctx)

	switch {
	case err != nil && l.LogLevel >= olog.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		logger.Error(sql, zap.Error(err), zap.Int64("rows", rows), zap.String("elapsed", elapsed.String()))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= olog.Warn:
		sql, rows := fc()

		logger.Warn(sql, zap.Int64("rows", rows), zap.String("elapsed", elapsed.String()))
	case l.LogLevel >= olog.Info:
		sql, rows := fc()
		logger.Debug(sql, zap.Int64("rows", rows), zap.String("elapsed", elapsed.String()))
	}
}

func (l GormLogger) logger(ctx context.Context) *zap.Logger {
	logger := handler.logger.WithOptions(zap.AddCallerSkip(1))
	if l.Context != nil {
		logger = logger.With(l.Context(ctx)...)
	}
	if trace := handler.trace(ctx); trace != nil {
		logger = logger.With(*trace)
	}
	return logger
}
