package logger

import (
	"context"
	"goapi/conf"
	"goapi/common/trace"
	"io"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/atopx/clever"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var handler *handle

type handle struct {
	logger   *zap.Logger
}

func (h *handle) trace(ctx context.Context) *zapcore.Field {
	if traceId, ok := trace.GetTraceID(ctx); ok {
		field := zap.String(string(trace.ContextTrace), traceId)
		return &field
	}
	return nil
}

func (h *handle) output(ctx context.Context, level zapcore.Level, message string, fields ...zapcore.Field) {
	if entity := handler.logger.Check(level, message); entity != nil {
		if trace := h.trace(ctx); trace != nil {
			fields = append(fields, *trace)
		}
		entity.Write(fields...)
	}
}

func writer(cfg *conf.LoggerConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.Filepath,
		MaxSize:    cfg.Maxsize,
		MaxAge:     cfg.Maxage,
		MaxBackups: cfg.Backups,
		LocalTime:  true,
	}
}

func Setup(cfg *conf.LoggerConfig) error {
	var loggerLevel = new(zapcore.Level)
	if err := loggerLevel.UnmarshalText(clever.Bytes(cfg.Level)); err != nil {
		return err
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}), zapcore.AddSync(writer(cfg)), loggerLevel)
	handler = &handle{
		logger:   zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)),
	}
	return nil
}

func Logger() *zap.Logger {
	return handler.logger
}

func Print(message string, fields ...zapcore.Field) {
	handler.logger.Info(message, fields...)
}

func Debug(ctx context.Context, message string, fields ...zapcore.Field) {
	handler.output(ctx, zap.DebugLevel, message, fields...)
}

func Info(ctx context.Context, message string, fields ...zapcore.Field) {
	handler.output(ctx, zap.InfoLevel, message, fields...)
}

func Warn(ctx context.Context, message string, fields ...zapcore.Field) {
	handler.output(ctx, zap.WarnLevel, message, fields...)
}

func Error(ctx context.Context, message string, fields ...zapcore.Field) {
	handler.output(ctx, zap.ErrorLevel, message, fields...)
}

func Fatal(ctx context.Context, message string, fields ...zapcore.Field) {
	handler.output(ctx, zap.FatalLevel, message, fields...)
}

func Panic(ctx context.Context, message string, fields ...zapcore.Field) {
	handler.output(ctx, zap.PanicLevel, message, fields...)
}
