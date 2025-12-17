package tests

import (
	"goapi/conf"
	"os"

	"github.com/spf13/viper"
)

func init() {
	// use conf/config.test.yaml
	os.Chdir("../")
	viper.SetConfigName("config.test")
	if _, err := conf.Load(); err != nil {
		panic(err)
	}
}
