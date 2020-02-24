package logger

import (
	"github.com/Lefthander/otus-go-antibruteforce/config"
	"go.uber.org/zap"
)

// GetLogger returns a zap logger in accordance with configuration settings
func GetLogger(cfg *config.ServiceConfig) (*zap.Logger, error) {
	var (
		l   *zap.Logger
		err error
	)

	switch cfg.LogMode {
	case "debug":
		l = zap.NewNop()
		return l, nil
	case "prod":
		l, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	case "dev":
		l, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	default:
		l = zap.NewExample()
	}
	return l, nil
}
