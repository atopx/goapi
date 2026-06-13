package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goapi/conf"
	"goapi/internal/common/handle"
	"goapi/internal/common/logger"
	"goapi/internal/model"
	"goapi/internal/scheduler"
	"goapi/internal/server"
	"goapi/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := conf.Load(); err != nil {
		log.Fatalf("load config: %v", err)
	}
	cfg := conf.Get()

	if err := logger.Setup(cfg.Logger); err != nil {
		log.Fatalf("setup logger: %v", err)
	}

	bootCtx := context.Background()
	logger.Info(bootCtx, "starting app",
		slog.String("name", cfg.AppName),
		slog.String("version", cfg.AppVersion),
	)

	switch cfg.Mode {
	case gin.DebugMode, gin.TestMode, gin.ReleaseMode:
		gin.SetMode(cfg.Mode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := pkg.NewDBClient(cfg.Database, logger.DbLogger())
	if err != nil {
		log.Fatalf("setup db: %v", err)
	}
	handle.Db = db
	if err := handle.Db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}
	handle.Redis = pkg.NewRedisClient(cfg.Redis)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	sch, err := scheduler.New(ctx, cfg.Scheduler)
	if err != nil {
		log.Fatalf("init scheduler: %v", err)
	}
	sch.Start()

	srv := server.New(cfg.Server)
	serverErr := make(chan error, 1)
	go func() {
		logger.Info(bootCtx, "http listening", slog.String("addr", cfg.Server.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		logger.Error(bootCtx, "server error", slog.Any("error", err))
		os.Exit(1)
	case <-ctx.Done():
		logger.Info(bootCtx, "shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "http shutdown failed", slog.Any("error", err))
	}
	if err := sch.Shutdown(); err != nil {
		logger.Error(shutdownCtx, "scheduler shutdown failed", slog.Any("error", err))
	}
	if sqlDB, err := handle.Db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	_ = handle.Redis.Close()
	logger.Info(bootCtx, "shutdown complete")
}
