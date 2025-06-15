package nightscout

import (
	"gabe565.com/trmnl-nightscout/internal/config"
)

type Properties struct {
	Bgnow     Reading   `json:"bgnow"`
	Delta     Delta     `json:"delta"`
	Direction Direction `json:"direction"`
}

func (p Properties) String(conf *config.Config) string {
	result := p.Bgnow.DisplayBg(conf.Units) +
		" " + p.Bgnow.Arrow()
	if delta := p.Delta.Display(conf.Units); delta != "" {
		result += " " + p.Delta.Display(conf.Units)
	}
	if rel := p.Bgnow.Mills.Relative(true); rel != "" {
		result += " [" + p.Bgnow.Mills.Relative(true) + "]"
	}
	return result
}
