package imaging

import (
	"image"
	"image/color"
)

type DotImage struct {
	gap     image.Point
	XOffset int
	palette color.Palette
}

func NewDots(gap image.Point, stagger bool) *DotImage {
	return (&DotImage{
		palette: Palette1Bit(),
	}).SetGap(gap, stagger)
}

func (d *DotImage) ColorModel() color.Model { return d.palette }
func (d *DotImage) Bounds() image.Rectangle { return image.Rect(-1e9, -1e9, 1e9, 1e9) }

func (d *DotImage) SetGap(gap image.Point, stagger bool) *DotImage {
	d.gap = gap.Add(image.Pt(1, 1))
	d.XOffset = 0
	if stagger {
		d.XOffset = d.gap.X / 2
	}
	return d
}

func (d *DotImage) SetForeground(c color.Color) *DotImage {
	d.palette[0] = c
	return d
}

func (d *DotImage) SetBackground(c color.Color) *DotImage {
	d.palette[1] = c
	return d
}

func (d *DotImage) At(x, y int) color.Color {
	if y%d.gap.Y == 0 {
		var shift int
		if (y/d.gap.Y)%2 == 1 {
			shift = d.XOffset
		}

		if (x-shift)%d.gap.X == 0 {
			return d.palette[0]
		}
	}
	return d.palette[1]
}
