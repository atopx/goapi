package conf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	AppName    string          `toml:"app_name"`
	AppVersion string          `toml:"app_version"`
	Server     *ServerConfig   `toml:"server"`
	Logger     *LoggerConfig   `toml:"logger"`
	Database   *DatabaseConfig `toml:"database"`
	Redis      *RedisConfig    `toml:"redis"`
	Scheduler  []WorkerConfig  `toml:"scheduler"`
}

type ServerConfig struct {
	Addr            string `toml:"addr"`
	ReadTimeout     int64  `toml:"read_timeout"`
	WriteTimeout    int64  `toml:"write_timeout"`
	MaxHeaderBytes  int    `toml:"max_header_bytes"`
	ShutdownTimeout int64  `toml:"shutdown_timeout"`
}

type LoggerConfig struct {
	Level     string `toml:"level"`
	Filepath  string `toml:"filepath"`
	Maxage    int    `toml:"maxage"`
	Maxsize   int    `toml:"maxsize"`
	Backups   int    `toml:"backups"`
	AddSource bool   `toml:"add_source"`
}

type WorkerConfig struct {
	Name    string         `toml:"name"`
	Spec    string         `toml:"spec"`
	Disable bool           `toml:"disable"`
	Args    map[string]any `toml:"args,omitempty"`
}

type DatabaseConfig struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	Name        string `toml:"name"`
	User        string `toml:"user"`
	Password    string `toml:"password"`
	SSLMode     string `toml:"ssl_mode"`
	MaxIdleConn int    `toml:"max_idle_conn"`
	MaxOpenConn int    `toml:"max_open_conn"`
	MaxIdleTime int64  `toml:"max_idle_time"`
	MaxLifeTime int64  `toml:"max_life_time"`
}

type RedisConfig struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	Password    string `toml:"password"`
	DB          int    `toml:"db"`
	PoolSize    int    `toml:"pool_size"`
	MaxLifeTime int64  `toml:"max_life_time"`
}

var config *Config

var searchPaths = []string{
	"config.toml",
	"conf/config.toml",
	"config/config.toml",
}

// Load 按 searchPaths 顺序查找 config.toml 并解析。
func Load() error {
	for _, p := range searchPaths {
		if _, err := os.Stat(p); err == nil {
			return load(p)
		}
	}
	return fmt.Errorf("config.toml not found in %v: %w", searchPaths, os.ErrNotExist)
}

func load(path string) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read config %s: %w", path, err)
	}
	cfg := new(Config)
	if err := toml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse toml %s: %w", path, err)
	}
	if cfg.Database == nil {
		return errors.New("config.database missing")
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}
	if cfg.Server != nil && cfg.Server.ShutdownTimeout <= 0 {
		cfg.Server.ShutdownTimeout = 10
	}
	config = cfg
	return nil
}

func Get() *Config {
	return config
}
