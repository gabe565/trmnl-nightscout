package config

import (
	"gabe565.com/trmnl-nightscout/internal/util"
	"github.com/caarlos0/env/v11"
)

func Load() (*Config, error) {
	conf, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	conf.Version = util.GetCommit()
	return &conf, nil
}
