package trmnl

import (
	"encoding/json"
	"time"
)

type Redirect struct {
	Filename    string        `json:"filename"`
	URL         string        `json:"url"`
	RefreshRate time.Duration `json:"refresh_rate"`
}

func (r Redirect) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"filename":     r.Filename,
		"url":          r.URL,
		"refresh_rate": int(r.RefreshRate.Seconds()),
	}
	return json.Marshal(m)
}
