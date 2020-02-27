package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	// DefaultClientTimeOut is timeout for connection to the abf-srv
	DefaultClientTimeOut = 30 * time.Second
)

// ClientConfig contains paramters related to the client
type ClientConfig struct {
	Host              string
	Port              string
	ConnectionTimeOut time.Duration
}

func newClientCfg() *ClientConfig {
	return &ClientConfig{
		Host:              viper.GetString("abf-ctl-host"),
		Port:              viper.GetString("abf-ctl-port"),
		ConnectionTimeOut: viper.GetDuration("abf-ctl-timeout"),
	}
}

// GetClientCfg returns a configuration parameters related to the client
func GetClientCfg() *ClientConfig {
	viper.SetDefault("abf-ctl-host", "localhost")
	viper.SetDefault("abf-ctl-port", "9000")
	viper.SetDefault("abf-ctl-timeout", DefaultClientTimeOut)

	return newClientCfg()
}
