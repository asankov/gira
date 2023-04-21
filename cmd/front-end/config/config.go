package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port          int    `required:"true" split_words:"true"`
	LogLevel      string `default:"info" split_words:"true"`
	APIAddress    string `required:"true" split_words:"true"`
	SessionSecret string `required:"true" split_words:"true"`
	EnforceHTTPS  bool   `required:"true" split_words:"true"`
}

func NewFromEnv() (*Config, error) {
	var config Config
	if err := envconfig.Process("GIRA", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
