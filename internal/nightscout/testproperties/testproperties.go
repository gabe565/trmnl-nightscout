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
	Etag       = `W/"20-8b9f9edb2e2b1a9f5a8ffbf92a1a1c42f170a654"`
)

func init() { //nolint:gochecknoinits
	err := json.Unmarshal(JSON, &Properties)
	if err != nil {
		panic(err)
	}
}
