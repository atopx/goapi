package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"goapi/conf"
)

const TraceKey = "trace_id"

var (
	defaultLogger *slog.Logger
	currentLevel  = new(slog.LevelVar)
)

func parseLevel(s string) (slog.Level, error) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(s)); err != nil {
		return lvl, fmt.Errorf("parse log level %q: %w", s, err)
	}
	return lvl, nil
}

func Setup(cfg *conf.LoggerConfig) error {
	lvl, err := parseLevel(cfg.Level)
	if err != nil {
		return err
	}
	currentLevel.Set(lvl)

	defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(defaultLogger)
	return nil
}

// Default 返回包级 slog.Logger，供少数需要直接传递 logger 的场景使用。
func Default() *slog.Logger {
	if defaultLogger == nil {
		return slog.Default()
	}
	return defaultLogger
}

// Level 返回当前生效的日志等级。
func Level() slog.Level {
	return currentLevel.Level()
}

func log(ctx context.Context, lvl slog.Level, msg string, attrs ...slog.Attr) {
	l := Default()
	if !l.Enabled(ctx, lvl) {
		return
	}
	if ctx != nil {
		if v, ok := ctx.Value(TraceKey).(string); ok && v != "" {
			attrs = append(attrs, slog.String(TraceKey, v))
		}
	}
	l.LogAttrs(ctx, lvl, msg, attrs...)
}

func Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	log(ctx, slog.LevelDebug, msg, attrs...)
}

func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	log(ctx, slog.LevelInfo, msg, attrs...)
}

func Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	log(ctx, slog.LevelWarn, msg, attrs...)
}

func Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	log(ctx, slog.LevelError, msg, attrs...)
}

// Fatal 记录 Error 等级日志后退出进程，仅供初始化阶段使用。
func Fatal(ctx context.Context, msg string, attrs ...slog.Attr) {
	log(ctx, slog.LevelError, msg, attrs...)
	os.Exit(1)
}
