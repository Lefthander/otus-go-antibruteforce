package config

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ServiceConfig is a struct to store all configuration parameters of service
type ServiceConfig struct {
	ConstraintN       uint32
	ConstraintM       uint32
	ConstraintK       uint32
	TimeOut           uint32
	DbHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	BucketIdleTimeOut time.Duration
	BucketCapacity    uint32
}

// GetConfig reads configuration from the provided file
func (c *ServiceConfig) GetConfig(cfgname string) (*ServiceConfig, error) {
	viper.SetConfigName(strings.Split(filepath.Base(cfgname), ".")[0])
	viper.AddConfigPath(cfgname)

	err := viper.ReadInConfig()

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
