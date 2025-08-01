package nightscout

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gabe565.com/trmnl-nightscout/internal/bg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReading_Arrow(t *testing.T) {
	t.Parallel()

	type fields struct {
		Mean      json.Number
		Last      bg.BG
		Mills     Mills
		Index     json.Number
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"TripleUp", fields{Sgvs: []SGV{{Direction: "TripleUp"}}}, "↑↑↑"},
		{"DoubleUp", fields{Sgvs: []SGV{{Direction: "DoubleUp"}}}, "↑↑"},
		{"SingleUp", fields{Sgvs: []SGV{{Direction: "SingleUp"}}}, "↑"},
		{"FortyFiveUp", fields{Sgvs: []SGV{{Direction: "FortyFiveUp"}}}, "↗"},
		{"Flat", fields{Sgvs: []SGV{{Direction: "Flat"}}}, "→"},
		{"FortyFiveDown", fields{Sgvs: []SGV{{Direction: "FortyFiveDown"}}}, "↘"},
		{"SingleDown", fields{Sgvs: []SGV{{Direction: "SingleDown"}}}, "↓"},
		{"DoubleDown", fields{Sgvs: []SGV{{Direction: "DoubleDown"}}}, "↓↓"},
		{"TripleDown", fields{Sgvs: []SGV{{Direction: "TripleDown"}}}, "↓↓↓"},
		{"unknown", fields{}, "-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			assert.Equal(t, tt.want, r.Arrow())
		})
	}
}

func TestReading_String(t *testing.T) {
	t.Parallel()
	type fields struct {
		Mean      json.Number
		Last      bg.BG
		Mills     Mills
		Index     json.Number
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	type args struct {
		unit bg.Unit
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"simple",
			fields{
				Last:  100,
				Mills: Mills{time.Now()},
				Sgvs:  []SGV{{Direction: "Flat"}},
			},
			args{},
			"100 → [0m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			assert.Equal(t, tt.want, r.String(tt.args.unit))
		})
	}
}

func TestReading_DisplayBg(t *testing.T) {
	t.Parallel()
	type fields struct {
		Mean      json.Number
		Last      bg.BG
		Mills     Mills
		Index     json.Number
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	type args struct {
		units bg.Unit
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   string
	}{
		{"95", args{bg.Mgdl}, fields{Last: 95}, "95"},
		{"LOW", args{bg.Mgdl}, fields{Last: 39}, "LOW"},
		{"HIGH", args{bg.Mgdl}, fields{Last: 401}, "HIGH"},
		{"mmol", args{bg.Mmol}, fields{Last: 100}, "5.6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			assert.Equal(t, tt.want, r.DisplayBg(tt.args.units))
		})
	}
}

var normalReading = `{
  "mean": 100,
  "last": 100,
  "mills": %d,
  "sgvs": [
    {
      "_id": "a",
      "mgdl": 100,
      "mills": %d,
      "device": "xDrip-DexcomG5",
      "direction": "Flat",
      "filtered": 0,
      "unfiltered": 0,
      "noise": 1,
      "rssi": 100,
      "type": "sgv",
      "scaled": 100
    }
  ]
}`

var lowReading = `{
  "sgvs": [
    {
      "_id": "a",
      "mgdl": 39,
      "mills": %d,
      "device": "xDrip-DexcomG5",
      "direction": "Flat",
      "filtered": 0,
      "unfiltered": 0,
      "noise": 1,
      "rssi": 100,
      "type": "sgv",
      "scaled": 39
    }
  ]
}`

func TestReading_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	now := time.Now()

	type fields struct {
		Mean      json.Number
		Last      bg.BG
		Mills     Mills
		Index     json.Number
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			"simple",
			fields{},
			args{[]byte(fmt.Sprintf(normalReading, now.UnixMilli(), now.UnixMilli()))},
			require.NoError,
		},
		{
			"low",
			fields{},
			args{[]byte(fmt.Sprintf(lowReading, now.UnixMilli()))},
			require.NoError,
		},
		{
			"error",
			fields{},
			args{[]byte("{")},
			require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			tt.wantErr(t, r.UnmarshalJSON(tt.args.bytes))
		})
	}
}
