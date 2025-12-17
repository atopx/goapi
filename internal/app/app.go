package app

import (
	"context"
	"fmt"
	"sync"

	"goapi/common/logger"
	"goapi/conf"
	"goapi/internal/infra"

	"go.uber.org/zap"
)

var (
	infrastructure *infra.Infrastructure
	once           sync.Once
)

func Initialize(ctx context.Context) {
	once.Do(func() {
		cfg, err := conf.Load()
		if err != nil {
			panic(fmt.Errorf("load config failed: %w", err))
		}
		if err = logger.Setup(cfg.Logger); err != nil {
			panic(fmt.Errorf("initialize logger failed: %w", err))
		}
		i, err := infra.NewInfrastructure(cfg)
		if err != nil {
			panic(fmt.Errorf("initialize infrastructure failed: %w", err))
		}
		infrastructure = i
		logger.Info(ctx, "application initialized", zap.String("app", cfg.AppName), zap.String("version", cfg.AppVersion))
	})
}

func Infra() *infra.Infrastructure {
	return infrastructure
}

func Close(ctx context.Context) {
	if infrastructure != nil {
		if err := infrastructure.Close(); err != nil {
			logger.Fatal(ctx, "infrastructure close failed")
		}
	}
}
