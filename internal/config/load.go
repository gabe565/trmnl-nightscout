package config

import (
	"github.com/caarlos0/env/v11"
)

func Load() (*Config, error) {
	conf, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
