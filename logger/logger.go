package logger

import (
	"github.com/Lefthander/otus-go-antibruteforce/config"
	"go.uber.org/zap"
)

// GetLogger returns a zap logger in accordance with configuration settings
func GetLogger(cfg *config.LoggerConfig) (*zap.SugaredLogger, error) {
	var (
		l   *zap.Logger
		err error
	)

	switch cfg.Environment {
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

	// err = l.Sync() I'm not sure that it's needed.

	if err != nil {
		return nil, err
	}

	return l.Sugar(), nil
}
