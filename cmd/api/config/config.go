package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port     int       `required:"true" split_words:"true"`
	Secret   string    `required:"true" split_words:"true"`
	UseSSL   bool      `required:"true" split_words:"true"`
	LogLevel string    `required:"true" split_words:"true"`
	DB       *DBConfig `required:"true" split_words:"true"`
}

type DBConfig struct {
	Host     string `required:"true" split_words:"true"`
	Port     int    `required:"true" split_words:"true"`
	Password string `required:"true" split_words:"true"`
	User     string `required:"true" split_words:"true"`
	Name     string `required:"true" split_words:"true"`
}

func NewFromEnv() (*Config, error) {
	var config Config
	if err := envconfig.Process("GIRA", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
