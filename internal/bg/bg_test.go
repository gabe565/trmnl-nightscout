package bg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBG_ToMmol(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		m    BG
		want float64
	}{
		{"100", BG(100), 5.55},
		{"50", BG(50), 2.775},
		{"300", BG(300), 16.65},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.InDelta(t, tt.want, tt.m.Mmol(), 0.001)
		})
	}
}
