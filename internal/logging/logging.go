package logging

import (
	"go-cloud-camp/internal/config"
	"sync"

	"go.uber.org/zap"
)

var instance *Logger
var once sync.Once

type Logger struct {
	*zap.SugaredLogger
}

// GetLogger function
func GetLogger(cfg config.LoggingParams) (*Logger, error) {
	var err error
	var logger *zap.Logger

	once.Do(func() {
		if *cfg.IsDebug {
			logger, err = zap.NewDevelopment()
		} else {
			logger, err = zap.NewProduction()
		}

		if err == nil {
			instance = &Logger{logger.Sugar()}
		}
	})

	return instance, err
}
