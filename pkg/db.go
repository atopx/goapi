package pkg

import (
	"fmt"
	"goapi/conf"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	olog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// NewDBClient 创建一个新的数据库客户端, 如果没有多类型的数据库，可以删除这个函数
func NewDBClient(cfg *conf.DatabaseConfig, logger olog.Interface) (*gorm.DB, error) {
	switch cfg.Type {
	case "mysql":
		return NewMySQLClient(cfg, logger)
	case "pgsql":
		return NewPgSQLClient(cfg, logger)
	case "sqlite":
		return NewSqliteClient(cfg, logger)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}

func NewMySQLClient(cfg *conf.DatabaseConfig, logger olog.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)
	return dbClient(mysql.Open(dsn), cfg, logger)
}

func NewPgSQLClient(cfg *conf.DatabaseConfig, logger olog.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Password)
	return dbClient(postgres.Open(dsn), cfg, logger)
}

func NewSqliteClient(cfg *conf.DatabaseConfig, logger olog.Interface) (*gorm.DB, error) {
	return dbClient(sqlite.Open(cfg.Name), cfg, logger)
}

func dbClient(dialector gorm.Dialector, cfg *conf.DatabaseConfig, logger olog.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: false,
		NowFunc:                func() time.Time { return time.Now().Local() },
		PrepareStmt:            true,
		CreateBatchSize:        1000,
		Logger:                 logger,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)
	sqlDb.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)
	return db, nil
}
