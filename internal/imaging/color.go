package imaging

import "image/color"

func Palette1Bit() color.Palette {
	return color.Palette{color.Black, color.White}
}
