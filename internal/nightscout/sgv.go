package nightscout

import (
	"encoding/json"

	"gabe565.com/trmnl-nightscout/internal/bg"
)

type SGV struct {
	ID         string      `json:"_id"`
	Device     string      `json:"device"`
	Direction  string      `json:"direction"`
	Filtered   json.Number `json:"filtered"`
	Mgdl       bg.BG       `json:"mgdl"`
	Mills      Mills       `json:"mills"`
	Noise      json.Number `json:"noise"`
	RSSI       json.Number `json:"rssi"`
	Scaled     json.Number `json:"scaled"`
	Type       string      `json:"type"`
	Unfiltered json.Number `json:"unfiltered"`
}

type SGVv1 struct {
	ID         string      `json:"_id"`
	Device     string      `json:"device"`
	Date       Mills       `json:"date"`
	SGV        bg.BG       `json:"sgv"`
	Delta      float64     `json:"delta"`
	Direction  string      `json:"direction"`
	Type       string      `json:"type"`
	Filtered   json.Number `json:"filtered"`
	Unfiltered json.Number `json:"unfiltered"`
	RSSI       json.Number `json:"rssi"`
	Noise      json.Number `json:"noise"`
	Mills      Mills       `json:"mills"`
}
