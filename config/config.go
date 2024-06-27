package config

import "github.com/caarlos0/env/v11"

type Config struct {
	JWTKey string `env:"JWT_KEY"`
}

func LoadEnvConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
