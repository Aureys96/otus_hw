package logger

import (
	"log"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
)

func New(appConfig *config.Config) *zap.Logger {
	var cfg zap.Config

	if appConfig.Production {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	al := zap.NewAtomicLevel()
	if err := al.UnmarshalText([]byte(appConfig.Logger.Level)); err != nil {
		log.Fatalln("Error while unmarshalling config", err)
	}

	cfg.Level.SetLevel(al.Level())

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalln("Error while building logger", err)
	}
	defer logger.Sync()

	return logger
}
