package config

import (
	"log"

	"github.com/spf13/viper"
)

// GetConfig reads configuration from the provided file
func GetConfig() error {
	if viper.GetString("config") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("abf")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Warning: Config file not found. Using defaults...", err)
		} else {
			log.Println("Error during processing config", err)
			return err
		}
	}

	return nil
}
