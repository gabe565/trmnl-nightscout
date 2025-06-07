package config

import "time"

//go:generate go tool envdoc -output ../../config.md
type Config struct {
	Version string `toml:"-"`

	// HTTP server bind address.
	ListenAddress string `env:"LISTEN_ADDRESS,notEmpty" envDefault:":8080"`
	// Get client IP address from the "Real-IP" header.
	RealIPHeader bool `env:"REAL_IP_HEADER"          envDefault:"false"`
	// Token required to access the API. If set, the value must be provided as a `token` query parameter.
	AccessToken string `env:"ACCESS_TOKEN"`
	// This app's public URL.
	PublicURL string `env:"PUBLIC_URL,notEmpty"`
	// Nightscout base URL
	NightscoutURL string `env:"NIGHTSCOUT_URL,notEmpty"`
	// Nightscout token. Using an access token is recommended instead of the API secret.
	NightscoutToken string `env:"NIGHTSCOUT_TOKEN"`
	// Blood sugar unit. (one of: mg/dL, mmol/L)
	Units Unit `env:"NIGHTSCOUT_UNITS"`
	// Time to wait before the next reading should be ready.\nIn testing, this seems to be about 20s behind, so the default is 30s to be safe.\nYour results may vary.
	FetchDelay time.Duration `env:"FETCH_DELAY"             envDefault:"30s"`
	// Normally, readings will be fetched when ready (after ~5m).\nThis interval will be used if the next reading time cannot be estimated due to sensor warm-up, missed readings, errors, etc.
	FallbackInterval time.Duration `env:"FALLBACK_INTERVAL"       envDefault:"30s"`

	// How far back in time the graph should go.
	GraphDuration time.Duration `env:"GRAPH_DURATION" envDefault:"6h"`
	// Where to draw the upper line.
	HighThreshold float64 `env:"HIGH_THRESHOLD" envDefault:"200"`
	// Where to draw the lower line.
	LowThreshold float64 `env:"LOW_THRESHOLD"  envDefault:"70"`
}
