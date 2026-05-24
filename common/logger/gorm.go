package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	olog "gorm.io/gorm/logger"
)

type GormLogger struct {
	LogLevel                  olog.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// DbLogger 根据当前全局日志等级派生 GORM 日志等级。
func DbLogger() GormLogger {
	var level olog.LogLevel
	switch Level() {
	case slog.LevelDebug:
		level = olog.Info
	case slog.LevelInfo, slog.LevelWarn:
		level = olog.Warn
	case slog.LevelError:
		level = olog.Error
	default:
		level = olog.Silent
	}
	return GormLogger{
		LogLevel:                  level,
		SlowThreshold:             500 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
	}
}

func (l GormLogger) LogMode(level olog.LogLevel) olog.Interface {
	clone := l
	clone.LogLevel = level
	return clone
}

func (l GormLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.LogLevel < olog.Info {
		return
	}
	Debug(ctx, fmt.Sprintf(msg, args...))
}

func (l GormLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.LogLevel < olog.Warn {
		return
	}
	Warn(ctx, fmt.Sprintf(msg, args...))
}

func (l GormLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.LogLevel < olog.Error {
		return
	}
	Error(ctx, fmt.Sprintf(msg, args...))
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= olog.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= olog.Error &&
		(!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		Error(ctx, "gorm query failed",
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsed.String()),
			slog.Any("error", err),
		)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= olog.Warn:
		sql, rows := fc()
		Warn(ctx, "gorm slow query",
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsed.String()),
		)
	case l.LogLevel >= olog.Info:
		sql, rows := fc()
		Debug(ctx, "gorm query",
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.String("elapsed", elapsed.String()),
		)
	}
}
