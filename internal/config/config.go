package config

import (
	"github.com/TinyMarcus/avito-tech-task/internal/db"
	"github.com/TinyMarcus/avito-tech-task/internal/logger"
	"github.com/kelseyhightower/envconfig"
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
