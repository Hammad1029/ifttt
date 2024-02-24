package config

import (
	"generic/utils"

	"github.com/spf13/viper"
)

var config *viper.Viper
var schemas *viper.Viper

func viperInit() {
	config = viper.New()
	schemas = viper.New()
	readEnv(config, "env", "json", "./config")
	readEnv(schemas, "schemas", "json", "./config")
}

func readEnv(config *viper.Viper, fileName string, fileType string, location string) {
	*config = *(viper.New())
	(*config).SetConfigName(fileName)
	(*config).SetConfigType(fileType)
	(*config).AddConfigPath(location)
	err := (*config).ReadInConfig()
	utils.HandleError(err, "fatal error config file")
}

func GetConfig() *viper.Viper {
	return config
}

func GetConfigProp(key string) string {
	return config.GetString(key)
}

func GetSchemas() *viper.Viper {
	return schemas
}

func GetSchemasProp(key string) string {
	return schemas.GetString(key)
}
