package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	AppName    string        `toml:"app_name"`
	AppVersion string        `toml:"app_version"`
	Server     *ServerConfig `toml:"server"`
	Logger     *LoggerConfig `toml:"logger"`
}

type ServerConfig struct {
	Addr            string `toml:"addr"`
	ReadTimeout     int64  `toml:"read_timeout"`
	WriteTimeout    int64  `toml:"write_timeout"`
	MaxHeaderBytes  int    `toml:"max_header_bytes"`
	ShutdownTimeout int64  `toml:"shutdown_timeout"`
}

type LoggerConfig struct {
	Level string `toml:"level"`
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
	config = cfg
	return nil
}

func Get() *Config {
	return config
}
