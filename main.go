package main

import (
	"context"
	_ "embed"
	"goapi/common/logger"
	"goapi/internal/app"
	"goapi/internal/scheduler"
	"goapi/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

func main() {

	app.Initialize(context.Background())
	cfg := app.Infra().Config()
	log.Printf("start app %s, current version is %s", cfg.AppName, cfg.AppVersion)

	if logger.Logger().Level() == zapcore.DebugLevel {
		// todo redoc openapi
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	errors := make(chan error, 1)

	s := scheduler.New(context.Background(), cfg.Scheduler)
	s.Start()
	server.Start(cfg.Server, errors)

	runloop(errors)
}

func runloop(errors <-chan error) {
	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case err := <-errors:
		log.Printf("runloop service error: %s", err.Error())
		break
	case s := <-notify:
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break
		}
	}
	app.Close(context.Background())
}
