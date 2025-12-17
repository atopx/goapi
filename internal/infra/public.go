package infra

import (
	"context"
	"goapi/conf"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func (i *Infrastructure) DB(ctx context.Context) *gorm.DB {
	return i.db.WithContext(ctx)
}

func (i *Infrastructure) RDB() *redis.Client {
	return i.rdb
}

func (i *Infrastructure) Config() *conf.Config {
	return i.config
}
