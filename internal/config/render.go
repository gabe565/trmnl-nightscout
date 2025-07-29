package config

import (
	"net/url"
	"strings"
	"time"

	"gabe565.com/trmnl-nightscout/internal/bg"
	"github.com/go-viper/mapstructure/v2"
)

type Render struct {
	// Blood sugar unit. (one of: mg/dL, mmol/L)
	Unit bg.Unit `env:"UNIT"`

	// Customize the time format. Use `3:04 PM` for 12-hour time or `15:04` for 24-hour. See [time](https://pkg.go.dev/time) for more details.
	TimeFormat string `env:"TIME_FORMAT" envDefault:"3:04 PM"`

	// How far back in time the graph should go.
	GraphDuration time.Duration `env:"GRAPH_DURATION" envDefault:"6h"`
	// Minimum X-axis value.
	GraphMin int `env:"GRAPH_MIN"      envDefault:"40"`
	// Maximum X-axis value.
	GraphMax int `env:"GRAPH_MAX"      envDefault:"300"`

	// Where to draw the upper line.
	HighThreshold float64 `env:"HIGH_THRESHOLD"        envDefault:"200"`
	// Background shade above the high threshold line. Value must be between 0-255.
	HighBackgroundShade uint8 `env:"HIGH_BACKGROUND_SHADE" envDefault:"245"`

	// Where to draw the lower line.
	LowThreshold float64 `env:"LOW_THRESHOLD"        envDefault:"70"`
	// Background shade below the low threshold line. Value must be between 0-255.
	LowBackgroundShade uint8 `env:"LOW_BACKGROUND_SHADE" envDefault:"237"`

	// Render with a black background and a white foreground.
	Invert bool `env:"INVERT"`
	// Invert colors when below this value. (Stacks with the `INVERT` option)
	InvertBelow float64 `env:"INVERT_BELOW" envDefault:"55"`
	// Invert colors when above this value. (Stacks with the `INVERT` option)
	InvertAbove float64 `env:"INVERT_ABOVE" envDefault:"300"`

	// Output color mode. 2-bit will be antialiased and dithering will be higher quality, but requires TRMNL firmware v1.6.0+. (one of 1bit, 2bit)
	ColorMode ColorMode `env:"COLOR_MODE" envDefault:"1bit"`
}

func (r *Render) UnmarshalQuery(q url.Values) error {
	m := map[string]string{}
	for k, v := range q {
		if len(v) != 0 && k != "token" {
			m[strings.ToUpper(k)] = v[0]
		}
	}

	if len(m) != 0 {
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook:       mapstructure.TextUnmarshallerHookFunc(),
			WeaklyTypedInput: true,
			Result:           r,
			TagName:          "env",
		})
		if err != nil {
			return err
		}

		if err := decoder.Decode(m); err != nil {
			return err
		}
	}

	if r.Unit == bg.Mmol {
		if r.HighThreshold > 39 {
			r.HighThreshold = bg.BG(r.HighThreshold).Mmol()
		}
		if r.LowThreshold > 39 {
			r.LowThreshold = bg.BG(r.LowThreshold).Mmol()
		}
		if r.InvertAbove > 39 {
			r.InvertAbove = bg.BG(r.InvertAbove).Mmol()
		}
		if r.InvertBelow > 39 {
			r.InvertBelow = bg.BG(r.InvertBelow).Mmol()
		}
		if r.GraphMin > 39 {
			r.GraphMin = 2
		}
		if r.GraphMax > 39 {
			r.GraphMax = 16
		}
	}

	return nil
}
