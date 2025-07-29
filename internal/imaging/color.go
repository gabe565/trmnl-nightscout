package imaging

import "image/color"

func Palette1Bit() color.Palette {
	return color.Palette{color.Black, color.White}
}

func Palette2Bit() color.Palette {
	return color.Palette{color.Black, color.Gray{Y: 0x55}, color.Gray{Y: 0xAA}, color.White}
}
