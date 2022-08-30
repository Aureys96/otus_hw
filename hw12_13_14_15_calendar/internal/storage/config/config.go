package config

type DBConfig struct {
	SQL      bool   `koanf:"sql"`
	Dsn      string `koanf:"dsn"`
	Inmemory bool   `koanf:"inmemory"`
}
