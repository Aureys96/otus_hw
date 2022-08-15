package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
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
