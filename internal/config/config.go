package config

import (
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Host string `validate:"required,hostname|ip"`
	Port int    `validate:"required,min=1,max=65535"`
}

func LoadConfig(host string, port int) (*Config, error) {
	return &Config{
		Host: host,
		Port: port,
	}, nil
}

func (c *Config) Validate() error {
	validator := validator.New()
	return validator.Struct(c)
}
