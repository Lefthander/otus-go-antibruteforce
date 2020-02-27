package config

import "github.com/spf13/viper"

// DBConfig contains all parameters related to DB connection
type DBConfig struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func newDBCfg() *DBConfig {
	return &DBConfig{
		DBHost: viper.GetString("abf-db-host"),
		DBPort: viper.GetString("abf-db-port"),
		DBUser: viper.GetString("abf-db-username"),
		DBPass: viper.GetString("abf-db-password"),
		DBName: viper.GetString("abf-db-dbname"),
	}
}

// GetDBCfg returns a DB configuration
func GetDBCfg() *DBConfig {
	viper.SetDefault("abf-db-host", "localhost")
	viper.SetDefault("abf-db-port", "5432")
	viper.SetDefault("abf-db-username", "postgres")
	viper.SetDefault("abf-db-password", "")
	viper.SetDefault("abf-db-dbname", "postgres")
	return newDBCfg()
}
