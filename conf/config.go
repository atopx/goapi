package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   *ServerConfig            `yaml:"server"`
	Logger   *LoggerConfig            `yaml:"logger"`
	Database *DatabaseConfig          `yaml:"database"`
	Redis    *RedisConfig             `yaml:"redis"`
	Workers  map[string]*WorkerConfig `yaml:"workers"`
}

type ServerConfig struct {
	Addr           string `yaml:"addr"`
	ReadTimeout    int64  `yaml:"read_timeout"`
	WriteTimeout   int64  `yaml:"write_timeout"`
	MaxHeaderBytes int    `yaml:"max_header_bytes"`
}

type LoggerConfig struct {
	Level    string `yaml:"level"`
	Filepath string `yaml:"filepath"`
	Maxage   int    `yaml:"maxage"`
	Maxsize  int    `yaml:"maxsize"`
	Backups  int    `yaml:"backups"`
	Trace    string `yaml:"trace"`
}

type WorkerConfig struct {
	Spec   string `yaml:"spec"`
	Enable bool   `yaml:"enable"`
}

type DatabaseConfig struct {
	Type        string `yaml:"type"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Name        string `yaml:"name"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
	MaxOpenConn int    `yaml:"max_open_conn"`
	MaxIdleTime int64  `yaml:"max_idle_time"`
	MaxLifeTime int64  `yaml:"max_life_time"`
}

type RedisConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Password    string `yaml:"password"`
	DB          int    `yaml:"db"`
	PoolSize    int    `yaml:"pool_size"`
	MaxLifeTime int64  `yaml:"max_life_time"`
}

var config *Config

func Load() error {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	config = new(Config)
	return viper.Unmarshal(&config)
}

func Get() *Config {
	return config
}
