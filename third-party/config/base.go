package config

import (
	"log"

	"github.com/spf13/viper"
)

// Setup initialize configuration
var (
	config *Configuration
)

// Params = getConfig.Params
func Setup() *Configuration {
	var baseConfiguration *EnvModel
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&baseConfiguration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	config = baseConfiguration.UpdateConfiguration()

	log.Println("configurations loading successfully")
	return config
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return config
}
