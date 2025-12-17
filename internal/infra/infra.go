package infra

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"goapi/common/logger"
	"goapi/conf"
	"goapi/internal/model"
	"goapi/pkg"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(
		ProvideMySQL,
		ProvideRedis,
		ProvideInfrastructure,
	),
)

func ProvideMySQL(cfg *conf.Config) (*gorm.DB, error) {
	return pkg.NewDBClient(cfg.Database, logger.DbLogger())
}

func ProvideRedis(cfg *conf.Config) (*redis.Client, error) {
	return pkg.NewRedisClient(cfg.Redis), nil
}

type Infrastructure struct {
	db     *gorm.DB
	rdb    *redis.Client
	config *conf.Config
}

func ProvideInfrastructure(
	cfg *conf.Config,
	mysql *gorm.DB,
	redis *redis.Client,
) *Infrastructure {
	return &Infrastructure{
		db:     mysql,
		rdb:    redis,
		config: cfg,
	}
}

func NewInfrastructure(cfg *conf.Config) (*Infrastructure, error) {
	var infra *Infrastructure
	app := fx.New(
		fx.Supply(cfg),
		Module,
		fx.Populate(&infra),
		fx.NopLogger,
	)
	if err := app.Start(contextWithTrace()); err != nil {
		return nil, fmt.Errorf("failed to start infrastructure: %w", err)
	}
	// auto migrate based on config
	if cfg.Database.AutoMigrate {
		if err := infra.db.AutoMigrate(&model.User{}); err != nil {
			return nil, fmt.Errorf("auto migrate failed: %w", err)
		}
	}
	return infra, nil
}

func (i *Infrastructure) Close() error {
	var errs []error
	if i.db != nil {
		if sqlDB, err := i.db.DB(); err != nil {
			errs = append(errs, fmt.Errorf("db get error: %w", err))
		} else if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("db close error: %w", err))
		}
	}
	if i.rdb != nil {
		if err := i.rdb.Close(); err != nil {
			errs = append(errs, fmt.Errorf("redis close error: %w", err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}
	return nil
}

func contextWithTrace() context.Context {
	return context.Background()
}
