package config

import "github.com/spf13/viper"

type Config struct {
	Log Log `json:"log,omitempty"`
}

type Log struct {
	Level      string
	File       string
	Structured bool
}

func New() (*Config, error) {
	cfg := &Config{}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
