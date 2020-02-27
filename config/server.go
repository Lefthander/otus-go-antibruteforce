package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	// DefaultBucketIdleTimeOut is a time after the bucket will be removed in case of
	// it's idle (no requests to it)
	DefaultBucketIdleTimeOut = 5 * time.Minute

	// DefaultTimeOut is timeout for context
	DefaultTimeOut = 30 * time.Second
)

// ServiceConfig contains the parameters related to ABF Service
type ServiceConfig struct {
	ServiceHost       string
	ServicePort       string
	MonitorPort       string
	ConstraintN       uint32
	ConstraintM       uint32
	ConstraintK       uint32
	BucketIdleTimeOut time.Duration
	BucketCapacity    uint32
	TimeOut           time.Duration
}

func newServiceCfg() *ServiceConfig {
	return &ServiceConfig{
		ServiceHost:       viper.GetString("abf-srv-host"),
		ServicePort:       viper.GetString("abf-srv-port"),
		MonitorPort:       viper.GetString("abf-srv-mon-port"),
		ConstraintN:       viper.GetUint32("abf-srv-login-limit"),
		ConstraintM:       viper.GetUint32("abf-srv-password-limit"),
		ConstraintK:       viper.GetUint32("abf-srv-ipaddress-limit"),
		BucketIdleTimeOut: viper.GetDuration("abf-srv-bucket-idle-timeout"),
		BucketCapacity:    viper.GetUint32("abf-srv-bucket-capacity"),
		TimeOut:           viper.GetDuration("abf-srv-default-timeout"),
	}
}

// GetServiceCfg returns a configuration data for ABF Service
func GetServiceCfg() *ServiceConfig {
	viper.SetDefault("abf-srv-host", "localhost")
	viper.SetDefault("abf-srv-port", "8999")
	viper.SetDefault("abf-srv-login-limit", 10)
	viper.SetDefault("abf-srv-password-limit", 100)
	viper.SetDefault("abf-srv-ipaddress-limit", 1000)
	viper.SetDefault("abf-srv-bucket-idle-timeout", DefaultBucketIdleTimeOut)
	viper.SetDefault("abf-srv-bucket-capacity", 1)
	viper.SetDefault("abf-srv-default-timeout", DefaultTimeOut)

	return newServiceCfg()
}
