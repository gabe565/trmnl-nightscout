package nightscout

import (
	"testing"
	"time"

	"gabe565.com/trmnl-nightscout/internal/bg"
	"github.com/stretchr/testify/assert"
)

func TestProperties_String(t *testing.T) {
	t.Parallel()
	type fields struct {
		Bgnow     Reading
		Buckets   []Reading
		Delta     Delta
		Direction Direction
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
			"mgdl",
			fields{
				Bgnow: Reading{
					Last:  100,
					Mills: Mills{time.Now()},
					Sgvs:  []SGV{{Direction: "Flat"}},
				},
				Delta: Delta{Mgdl: 1},
			},
			args{},
			"100 → +1 [0m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := Properties{
				Bgnow: tt.fields.Bgnow,
				Delta: tt.fields.Delta,
			}
			assert.Equal(t, tt.want, p.String(tt.args.unit))
		})
	}
}
