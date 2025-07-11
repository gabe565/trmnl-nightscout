package nightscout

import (
	"encoding/json"
	"testing"

	"gabe565.com/trmnl-nightscout/internal/bg"
	"github.com/stretchr/testify/assert"
)

func TestDelta_Display(t *testing.T) {
	t.Parallel()
	type fields struct {
		Absolute     json.Number
		DisplayVal   string
		ElapsedMins  json.Number
		Interpolated bool
		Mean5MinsAgo json.Number
		Mgdl         bg.BG
		Previous     Reading
		Scaled       json.Number
		Times        Times
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
		{
			"mgdl",
			args{bg.Mgdl},
			fields{Mgdl: 1},
			"+1",
		},
		{
			"mmol",
			args{bg.Mmol},
			fields{Mgdl: 9},
			"+0.5",
		},
		{
			"mmol no decimal",
			args{bg.Mmol},
			fields{Mgdl: 0},
			"+0",
		},
		{
			"mmol negative",
			args{bg.Mmol},
			fields{Mgdl: -9},
			"-0.5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := Delta{
				Absolute:     tt.fields.Absolute,
				DisplayVal:   tt.fields.DisplayVal,
				ElapsedMins:  tt.fields.ElapsedMins,
				Interpolated: tt.fields.Interpolated,
				Mean5MinsAgo: tt.fields.Mean5MinsAgo,
				Mgdl:         tt.fields.Mgdl,
				Previous:     tt.fields.Previous,
				Scaled:       tt.fields.Scaled,
				Times:        tt.fields.Times,
			}
			assert.Equal(t, tt.want, d.Display(tt.args.units))
		})
	}
}
