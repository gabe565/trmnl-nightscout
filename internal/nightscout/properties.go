package nightscout

import (
	"gabe565.com/trmnl-nightscout/internal/bg"
)

type Properties struct {
	Bgnow     Reading   `json:"bgnow"`
	Delta     Delta     `json:"delta"`
	Direction Direction `json:"direction"`
}

func (p Properties) String(unit bg.Unit) string {
	result := p.Bgnow.DisplayBg(unit) +
		" " + p.Bgnow.Arrow()
	if delta := p.Delta.Display(unit); delta != "" {
		result += " " + p.Delta.Display(unit)
	}
	if rel := p.Bgnow.Mills.Relative(true); rel != "" {
		result += " [" + p.Bgnow.Mills.Relative(true) + "]"
	}
	return result
}
