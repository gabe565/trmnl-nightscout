package imaging

import (
	"image"
	"image/color"
)

type DotImage struct {
	spaceX, spaceY int
	XOffset        int
	palette        color.Palette
}

func NewDots(x, y int, stagger bool) *DotImage {
	return (&DotImage{
		palette: Palette1Bit(),
	}).SetSpacing(x, y, stagger)
}

func (d *DotImage) ColorModel() color.Model { return d.palette }
func (d *DotImage) Bounds() image.Rectangle { return image.Rect(-1e9, -1e9, 1e9, 1e9) }

func (d *DotImage) SetSpacing(x, y int, stagger bool) *DotImage {
	d.spaceX = x + 1
	d.spaceY = y + 1
	d.XOffset = 0
	if stagger {
		d.XOffset = d.spaceX / 2
	}
	return d
}

func (d *DotImage) At(x, y int) color.Color {
	if y%d.spaceY != 0 {
		return d.palette[1]
	}

	var shift int
	if (y/d.spaceY)%2 == 1 {
		shift = d.XOffset
	}

	if (x-shift)%d.spaceX == 0 {
		return d.palette[0]
	}
	return d.palette[1]
}
