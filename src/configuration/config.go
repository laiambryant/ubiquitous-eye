package configuration

import (
	env "github.com/goloop/env"
)

type Config struct {
	Level string `env:"level" def:"info"`
	File  string `env:"file"  def:"logs/app.log"`
}

func LoadConfig() (*Config, error) {
	err := env.Load("configuration/ue.config")
	return &Config{}, err
}
