package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"goapi/common/logger"
	"goapi/conf"
	"goapi/internal/server"

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

	if logger.Level() > slog.LevelDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

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
	case <-ctx.Done():
		logger.Info(bootCtx, "shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "http shutdown failed", slog.Any("error", err))
	}
	logger.Info(bootCtx, "shutdown complete")
}
