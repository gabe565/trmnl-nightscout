package config

import (
	"time"
)

//go:generate go tool envdoc -types Config -output ../../docs/envs.md

// Config is the base configuration. Note that many of these envs can be overridden per-request using [query parameters](query-params.md).
type Config struct {
	Version string `toml:"-"`

	// Nightscout base URL
	NightscoutURL string `env:"NIGHTSCOUT_URL,required"`

	// A URL that the TRMNL device can use to download the image from this app. It can be a public URL or an internal IP address as long as the TRMNL device is on the same network.
	ImageURL string `env:"IMAGE_URL,required"`

	// Nightscout token. Using an access token is recommended instead of the API secret.
	NightscoutToken string `env:"NIGHTSCOUT_TOKEN"`

	// HTTP server bind address.
	ListenAddress string `env:"LISTEN_ADDRESS" envDefault:":8080"`
	// Token required to access the API. If set, the value must be provided as a `token` query parameter.
	AccessToken string `env:"ACCESS_TOKEN"`
	// Get client IP address from the "Real-IP" header.
	RealIPHeader bool `env:"REAL_IP_HEADER" envDefault:"false"`

	// The interval that new readings are sent to Nightscout.
	UpdateInterval time.Duration `env:"UPDATE_INTERVAL"   envDefault:"5m"`
	// Time to wait before the next reading should be ready. In testing, this seems to be about 20s behind, so the default is 30s to be safe. Your results may vary.
	FetchDelay time.Duration `env:"FETCH_DELAY"       envDefault:"30s"`
	// Normally, readings will be fetched when ready (after ~5m). This interval will be used if the next reading time cannot be estimated due to sensor warm-up, missed readings, errors, etc.
	FallbackInterval time.Duration `env:"FALLBACK_INTERVAL" envDefault:"30s"`

	Render Render
}
