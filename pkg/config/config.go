package config

import (
	"os"

	"github.com/spf13/viper"
)

func Read() error {
	viper.AddConfigPath(os.Getenv("CONFIG_PATH_PREFIX") + "pkg/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
