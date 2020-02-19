package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// ServiceConfig is a struct to store all configuration parameters of service
type ServiceConfig struct {
	ConstraintN uint32
	ConstraintM uint32
	ConstraintK uint32
	DbHost      string
	DBPort      string
	DBUser      string
	DBPass      string
	DBName      string
}

// GetConfig reads configuration from the provided file
func (c *ServiceConfig) GetConfig(path string) (*ServiceConfig, error) {
	viper.SetConfigName(strings.Split(path.Base(path), ".")[0])
	viper.AddConfigPath(path)

	err := viper.ReadConfig()

	cfg := ServiceConfig{}

	if err != nil {

		log.Fatal("Failed to get the config file", err)
	}

	err = viper.Unmarshal(&cfg)

	if err != nil {

		log.Fatal("Failed to unmarshal config", err)
	}

	return &cfg, nil

}
