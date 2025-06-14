package imaging

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvert1Bit(t *testing.T) {
	img := image.NewPaletted(image.Rect(0, 0, 100, 100), Palette1Bit())
	InvertPaletted(img)
	for _, c := range img.Pix {
		assert.Equal(t, uint8(1), c)
	}
}
