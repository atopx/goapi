package main

import (
	_ "embed"
	"fmt"
	"goapi/common/handle"
	"goapi/common/logger"
	"goapi/conf"
	"goapi/internal/model"
	"goapi/internal/server"
	"goapi/internal/worker"
	"goapi/pkg"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

func main() {

	err := conf.Load()
	if err != nil {
		panic(fmt.Errorf("load config error: %s", err.Error()))
	}

	cfg := conf.Get()
	log.Printf("start app %s, current version is %s", cfg.AppName, cfg.AppVersion)

	if err = logger.Setup(cfg.Logger); err != nil {
		panic(fmt.Errorf("setup logger error: %s", err.Error()))
	}

	if logger.Logger().Level() == zapcore.DebugLevel {
		// todo redoc openapi
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	handle.Db, err = pkg.NewDBClient(cfg.Database, logger.DbLogger())
	if err != nil {
		panic(fmt.Errorf("setup db error: %s", err.Error()))
	}

	handle.Db.AutoMigrate(&model.User{})

	handle.Redis = pkg.NewRedisClient(cfg.Redis)

	errors := make(chan error, 1)

	worker.Start(cfg.Workers, errors)
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
}
