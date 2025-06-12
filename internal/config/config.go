package config

import "time"

//go:generate go tool envdoc -output ../../config.md
type Config struct {
	Version string `toml:"-"`

	// HTTP server bind address.
	ListenAddress string `env:"LISTEN_ADDRESS,notEmpty" envDefault:":8080"`
	// This app's public URL.
	PublicURL string `env:"PUBLIC_URL,notEmpty"`
	// Token required to access the API. If set, the value must be provided as a `token` query parameter.
	AccessToken string `env:"ACCESS_TOKEN"`
	// Get client IP address from the "Real-IP" header.
	RealIPHeader bool `env:"REAL_IP_HEADER"          envDefault:"false"`

	// Nightscout base URL
	NightscoutURL string `env:"NIGHTSCOUT_URL,notEmpty"`
	// Nightscout token. Using an access token is recommended instead of the API secret.
	NightscoutToken string `env:"NIGHTSCOUT_TOKEN"`

	// Blood sugar unit. (one of: mg/dL, mmol/L)
	Units Unit `env:"NIGHTSCOUT_UNITS"`
	// Customize the time format. Use `3:04 PM` for 12-hour time or `15:04` for 24-hour. See [time](https://pkg.go.dev/time) for more details.
	TimeFormat string `env:"TIME_FORMAT"      envDefault:"3:04 PM"`

	// How far back in time the graph should go.
	GraphDuration time.Duration `env:"GRAPH_DURATION" envDefault:"6h"`

	// Where to draw the upper line.
	HighThreshold float64 `env:"HIGH_THRESHOLD"        envDefault:"200"`
	// Background shade above the high threshold line. Value must be between 0-255.
	HighBackgroundShade uint8 `env:"HIGH_BACKGROUND_SHADE" envDefault:"250"`

	// Where to draw the lower line.
	LowThreshold float64 `env:"LOW_THRESHOLD"        envDefault:"70"`
	// Background shade below the low threshold line. Value must be between 0-255.
	LowBackgroundShade uint8 `env:"LOW_BACKGROUND_SHADE" envDefault:"247"`

	// Render with a black background and a white foreground.
	Invert bool `env:"INVERT"`
	// Invert colors when below this value. (Stacks with the `INVERT` option)
	InvertBelow int `env:"INVERT_BELOW" envDefault:"55"`
	// Invert colors when above this value. (Stacks with the `INVERT` option)
	InvertAbove int `env:"INVERT_ABOVE" envDefault:"300"`

	// Time to wait before the next reading should be ready. In testing, this seems to be about 20s behind, so the default is 30s to be safe. Your results may vary.
	FetchDelay time.Duration `env:"FETCH_DELAY"       envDefault:"30s"`
	// Normally, readings will be fetched when ready (after ~5m). This interval will be used if the next reading time cannot be estimated due to sensor warm-up, missed readings, errors, etc.
	FallbackInterval time.Duration `env:"FALLBACK_INTERVAL" envDefault:"30s"`
}
