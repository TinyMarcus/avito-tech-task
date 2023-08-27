package config

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/TinyMarcus/avito-tech-task/internal/db"
	"github.com/TinyMarcus/avito-tech-task/internal/logger"
)

type Config struct {
	Log  logger.LogConfig  `envconfig:"LOG"`
	Db   db.DatabaseConfig `envconfig:"DB"`
	Port string            `envconfig:"PORT"`
}

func New() (*Config, error) {
	var config Config

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
