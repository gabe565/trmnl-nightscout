package imaging

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvert1Bit(t *testing.T) {
	img := image.NewPaletted(image.Rect(0, 0, 100, 100), color.Palette{color.Black, color.White})
	InvertPaletted(img)
	for _, c := range img.Pix {
		assert.Equal(t, uint8(1), c)
	}
}
