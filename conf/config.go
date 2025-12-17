package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppName    string          `mapstruct:"app_name"`
	AppVersion string          `mapstruct:"app_version"`
	Server     *ServerConfig   `mapstruct:"server"`
	Logger     *LoggerConfig   `mapstruct:"logger"`
	Database   *DatabaseConfig `mapstruct:"database"`
	Redis      *RedisConfig    `mapstruct:"redis"`
	Scheduler  []WorkerConfig  `mapstruct:"scheduler"`
}

type ServerConfig struct {
	Addr           string `mapstruct:"addr"`
	ReadTimeout    int64  `mapstruct:"read_timeout"`
	WriteTimeout   int64  `mapstruct:"write_timeout"`
	MaxHeaderBytes int    `mapstruct:"max_header_bytes"`
}

type LoggerConfig struct {
	Level    string `mapstruct:"level"`
	Filepath string `mapstruct:"filepath"`
	Maxage   int    `mapstruct:"maxage"`
	Maxsize  int    `mapstruct:"maxsize"`
	Backups  int    `mapstruct:"backups"`
}

type WorkerConfig struct {
	Name    string         `mapstruct:"name"`
	Spec    string         `mapstruct:"spec"`
	Disable bool           `mapstruct:"disable"`
	Args    map[string]any `mapstruct:"args,omitempty"`
}

type DatabaseConfig struct {
	Type        string `mapstruct:"type"`
	Host        string `mapstruct:"host"`
	Port        int    `mapstruct:"port"`
	Name        string `mapstruct:"name"`
	User        string `mapstruct:"user"`
	Password    string `mapstruct:"password"`
	MaxIdleConn int    `mapstruct:"max_idle_conn"`
	MaxOpenConn int    `mapstruct:"max_open_conn"`
	MaxIdleTime int64  `mapstruct:"max_idle_time"`
	MaxLifeTime int64  `mapstruct:"max_life_time"`
	AutoMigrate bool   `mapstruct:"auto_migrate"`
}

type RedisConfig struct {
	Host        string `mapstruct:"host"`
	Port        int    `mapstruct:"port"`
	Password    string `mapstruct:"password"`
	DB          int    `mapstruct:"db"`
	PoolSize    int    `mapstruct:"pool_size"`
	MaxLifeTime int64  `mapstruct:"max_life_time"`
}

func Load() (*Config, error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("./config")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
