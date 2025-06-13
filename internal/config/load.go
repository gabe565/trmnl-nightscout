package config

import (
	"gabe565.com/trmnl-nightscout/internal/bg"
	"github.com/caarlos0/env/v11"
)

func Load() (*Config, error) {
	conf, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	if conf.Units == bg.Mmol {
		if conf.HighThreshold > 39 {
			conf.HighThreshold = bg.BG(conf.HighThreshold).Mmol()
		}
		if conf.LowThreshold > 39 {
			conf.LowThreshold = bg.BG(conf.LowThreshold).Mmol()
		}
		if conf.InvertAbove > 39 {
			conf.InvertAbove = bg.BG(conf.InvertAbove).Mmol()
		}
		if conf.InvertBelow > 39 {
			conf.InvertBelow = bg.BG(conf.InvertBelow).Mmol()
		}
	}

	return &conf, nil
}
