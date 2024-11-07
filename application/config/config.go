package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init() error {
	config = viper.New()
	config.SetConfigName("env")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	if err := config.ReadInConfig(); err != nil {
		return fmt.Errorf("could not read in config %s", err)
	}
	return nil
}

func GetConfig() *viper.Viper {
	return config
}

func GetConfigProp(key string) string {
	return config.GetString(key)
}
