package pkg

import (
	"fmt"
	"time"

	"goapi/conf"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	olog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDBClient(cfg *conf.DatabaseConfig, logger olog.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Password, cfg.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		NowFunc:                func() time.Time { return time.Now().Local() },
		PrepareStmt:            true,
		CreateBatchSize:        1000,
		Logger:                 logger,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)
	return db, nil
}
