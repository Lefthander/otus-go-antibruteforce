package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// GetConfig reads configuration from the provided file
func GetConfig(cfg string) error {
	viper.SetConfigName(filepath.Base(strings.Split(filepath.Base(cfg), ".")[0]))
	viper.AddConfigPath(filepath.Dir(cfg))
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
