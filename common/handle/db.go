package handle

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Db    *gorm.DB
	Redis *redis.Client
)

func DB(ctx context.Context) *gorm.DB {
	return Db.WithContext(ctx)
}

func RDB() *redis.Client {
	return Redis
}
