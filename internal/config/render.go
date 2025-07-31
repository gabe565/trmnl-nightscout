package config

import (
	"net/url"
	"strings"
	"time"

	"gabe565.com/trmnl-nightscout/internal/bg"
	"github.com/go-viper/mapstructure/v2"
	"gonum.org/v1/plot/vg"
)

type Render struct {
	// Blood sugar unit. (one of: mg/dL, mmol/L)
	Unit bg.Unit `env:"UNIT"`

	// Customize the time format. Use `3:04 PM` for 12-hour time or `15:04` for 24-hour. See [time](https://pkg.go.dev/time) for more details.
	TimeFormat string `env:"TIME_FORMAT" envDefault:"3:04 PM"`

	// How far back in time the graph should go.
	GraphDuration time.Duration `env:"GRAPH_DURATION" envDefault:"6h"`
	// Minimum X-axis value.
	GraphMin bg.BG `env:"GRAPH_MIN"      envDefault:"40"`
	// Maximum X-axis value.
	GraphMax bg.BG `env:"GRAPH_MAX"      envDefault:"300"`

	// Control the plot point stroke radius. Set to 0 to disable.
	PointStrokeRadius vg.Length `env:"POINT_STROKE_RADIUS" envDefault:"4"`

	// Where to draw the upper line.
	HighThreshold bg.BG `env:"HIGH_THRESHOLD" envDefault:"200"`
	// Where to draw the lower line.
	LowThreshold bg.BG `env:"LOW_THRESHOLD"  envDefault:"70"`

	// Render with a black background and a white foreground.
	Invert bool `env:"INVERT"`
	// Invert colors when below this value. (Stacks with the `INVERT` option)
	InvertBelow bg.BG `env:"INVERT_BELOW" envDefault:"55"`
	// Invert colors when above this value. (Stacks with the `INVERT` option)
	InvertAbove bg.BG `env:"INVERT_ABOVE" envDefault:"300"`

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
		if r.HighThreshold == 200 {
			r.HighThreshold = bg.NewMmol(11)
		} else if r.HighThreshold < 30 {
			r.HighThreshold = bg.NewMmol(r.HighThreshold)
		}
		if r.LowThreshold == 70 {
			r.LowThreshold = bg.NewMmol(4)
		} else if r.LowThreshold < 30 {
			r.LowThreshold = bg.NewMmol(r.LowThreshold)
		}
		if r.InvertAbove < 30 {
			r.InvertAbove = bg.NewMmol(r.InvertAbove)
		}
		if r.InvertBelow < 30 {
			r.InvertBelow = bg.NewMmol(r.InvertBelow)
		}
		if r.GraphMin == 40 {
			r.GraphMin = bg.NewMmol(2)
		} else if r.GraphMin < 30 {
			r.GraphMin = bg.NewMmol(r.GraphMin)
		}
		if r.GraphMax == 300 {
			r.GraphMax = bg.NewMmol(16)
		} else if r.GraphMax < 30 {
			r.GraphMax = bg.NewMmol(r.GraphMax)
		}
	}

	return nil
}
