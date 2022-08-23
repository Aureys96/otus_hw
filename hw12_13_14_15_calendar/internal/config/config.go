package config

import (
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/server/config"
	storageConfig "github.com/Aureys96/hw12_13_14_15_calendar/internal/storage/config"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Logger     LoggerConf             `koanf:"logger"`
	Server     config.ServerConfig    `koanf:"server"`
	DbConfig   storageConfig.DBConfig `koanf:"storage"`
	Production bool                   `koanf:"production"`
}

type LoggerConf struct {
	Level            string
	Encoding         string
	outputPaths      string
	ErrorOutputPaths []string
}

func NewConfig(configPath string) (*Config, error) {
	k := koanf.New(".")

	k.Load(confmap.Provider(map[string]interface{}{
		"logger.level":            "INFO",
		"logger.encoding":         "json",
		"logger.outputPaths":      []string{"stdout"},
		"logger.errorOutputPaths": []string{"stderr"},
		"server.host":             "localhost",
		"server.port":             "8080",
		"server.shutdownTime":     "5",
		"storage.inmemory":        true,
		"production":              false,
	}, "."), nil)
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		return nil, err
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
