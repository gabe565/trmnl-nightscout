package imaging

import "image/color"

//nolint:gochecknoglobals
var (
	Gray1 = color.Gray{Y: 0x55}
	Gray2 = color.Gray{Y: 0xAA}
)

func Palette1Bit() color.Palette {
	return color.Palette{color.Black, color.White}
}

func Palette2Bit() color.Palette {
	return color.Palette{color.Black, Gray1, Gray2, color.White}
}
