package config

import (
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/server/config"
	storageConfig "github.com/Aureys96/hw12_13_14_15_calendar/internal/storage/config"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf
	Server   config.ServerConfig `koanf:"server"`
	Inmemory bool
	DbConfig storageConfig.DBConfig
	// TODO
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
		"server.host":         "localhost",
		"server.port":         "8080",
		"server.shutdownTime": "5",
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

// TODO
