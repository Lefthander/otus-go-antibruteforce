package config

import "github.com/spf13/viper"

// LoggerConfig contains configuration of the Logger
type LoggerConfig struct {
	Verbose     bool
	Debug       bool
	Environment string
}

func newLoggerCfg() *LoggerConfig {
	return &LoggerConfig{
		Verbose:     viper.GetBool("abf-srv-logger-verbose"),
		Debug:       viper.GetBool("abf-srv-logger-debug"),
		Environment: viper.GetString("abf-srv-env"),
	}
}

// GetLoggerCfg returns the LoggerConfig within pre-defined settings in case of missing some...
func GetLoggerCfg() *LoggerConfig {
	viper.SetDefault("bf-srv-logger-verbose", false)
	viper.SetDefault("abf-srv-logger-debug", false)
	viper.SetDefault("abf-srv-env", "dev")
	return newLoggerCfg()
}
