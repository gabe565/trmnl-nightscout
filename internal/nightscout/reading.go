package nightscout

import (
	"encoding/json"
	"math"
	"strconv"

	"gabe565.com/trmnl-nightscout/internal/bg"
)

const (
	LowReading  = 39
	HighReading = 401
)

type Reading struct {
	Mean      json.Number `json:"mean"`
	Last      bg.BG       `json:"last"`
	Mills     Mills       `json:"mills"`
	Index     json.Number `json:"index,omitempty"`
	FromMills Mills       `json:"fromMills,omitempty"`
	ToMills   Mills       `json:"toMills,omitempty"`
	Sgvs      []SGV       `json:"sgvs"`
}

func (r *Reading) Arrow() string {
	var direction string
	if len(r.Sgvs) > 0 {
		direction = r.Sgvs[0].Direction
	}
	switch direction {
	case "TripleUp":
		return "↑↑↑"
	case "DoubleUp":
		return "↑↑"
	case "SingleUp":
		return "↑"
	case "FortyFiveUp":
		return "↗"
	case "Flat":
		return "→"
	case "FortyFiveDown":
		return "↘"
	case "SingleDown":
		return "↓"
	case "DoubleDown":
		return "↓↓"
	case "TripleDown":
		return "↓↓↓"
	default:
		return "-"
	}
}

func (r *Reading) String(unit bg.Unit) string {
	if r.Last == 0 {
		return ""
	}

	result := r.DisplayBg(unit) +
		" " + r.Arrow()
	if rel := r.Mills.Relative(true); rel != "" {
		result += " [" + r.Mills.Relative(true) + "]"
	}
	return result
}

func (r *Reading) UnmarshalJSON(bytes []byte) error {
	type rawReading Reading
	if err := json.Unmarshal(bytes, (*rawReading)(r)); err != nil {
		return err
	}

	// Last is unset if reading is out of range.
	// Will be set from sgvs.
	if r.Last == 0 && len(r.Sgvs) > 0 {
		r.Last = r.Sgvs[0].Mgdl
		r.Mills = r.Sgvs[0].Mills
	}

	return nil
}

func (r *Reading) DisplayBg(units bg.Unit) string {
	switch r.Last {
	case LowReading:
		return "LOW"
	case HighReading:
		return "HIGH"
	}

	if units == bg.Mmol {
		mmol := r.Last.Mmol()
		mmol = math.Round(mmol*10) / 10
		return strconv.FormatFloat(mmol, 'f', 1, 64)
	}

	return strconv.FormatFloat(r.Last.Mgdl(), 'f', -1, 64)
}
