package imaging

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvert1Bit(t *testing.T) {
	t.Run("fails when not 1-bit", func(t *testing.T) {
		img := image.NewPaletted(
			image.Rect(0, 0, 100, 100),
			color.Palette{color.Black, color.White, color.Gray{Y: 0x80}},
		)
		require.Error(t, Invert1Bit(img))
		for _, c := range img.Pix {
			assert.Equal(t, uint8(0), c)
		}
	})

	t.Run("works when 1-bit", func(t *testing.T) {
		img := image.NewPaletted(image.Rect(0, 0, 100, 100), color.Palette{color.Black, color.White})
		require.NoError(t, Invert1Bit(img))
		for _, c := range img.Pix {
			assert.Equal(t, uint8(1), c)
		}
	})
}
