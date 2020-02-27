package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// GetConfig reads configuration from the provided file
func GetConfig(cfgname string) error {
	viper.SetConfigName(strings.Split(filepath.Base(cfgname), ".")[0])
	viper.AddConfigPath(cfgname)

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
