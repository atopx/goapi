package pkg

import (
	"fmt"
	"goapi/conf"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *conf.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		ConnMaxLifetime: time.Duration(cfg.MaxLifeTime) * time.Second,
	})
}
