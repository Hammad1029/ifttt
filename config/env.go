package config

import "github.com/spf13/viper"

var config *viper.Viper

func readEnv() {
	config = viper.New()
	config.SetConfigName("env")
	config.SetConfigType("json")
	config.AddConfigPath("./config")
	err := config.ReadInConfig()
	handleError(err, "fatal error config file")
}

func GetConfig() *viper.Viper {
	return config
}

func GetConfigProp(key string) string {
	return config.GetString(key)
}
