package config

import "github.com/caarlos0/env/v11"

type Config struct {
	DBUrl string `env:"DB_URL"`
	Host  string `env:"SERVER_HOST" envDefault:":8080"`
}

// New initialize the project configuration.
func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
