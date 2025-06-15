//nolint:gochecknoglobals,gosmopolitan
package testproperties

import (
	_ "embed"
	"encoding/json"

	"gabe565.com/trmnl-nightscout/internal/nightscout"
)

var (
	//go:embed fetch_test_properties.json
	JSON       []byte
	Properties nightscout.Properties
	Etag       = `W/"20-4353d9446a4377f8bc10e267ff9be8c81ca4279a"`
)

func init() { //nolint:gochecknoinits
	err := json.Unmarshal(JSON, &Properties)
	if err != nil {
		panic(err)
	}
}
