package logger

import (
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
)

func New(appConfig *config.Config) *zap.Logger {
	//var cfg zap.Config
	cfg := zap.NewDevelopmentConfig()

	al := zap.NewAtomicLevel()
	err := al.UnmarshalText([]byte(appConfig.Logger.Level))
	cfg.Level.SetLevel(al.Level())

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}
