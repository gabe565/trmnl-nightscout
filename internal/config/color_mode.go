package config

import (
	"image/color"

	"gabe565.com/trmnl-nightscout/internal/imaging"
)

//go:generate go tool enumer -type ColorMode -trimprefix ColorMode -transform lower -text

//nolint:recvcheck
type ColorMode uint8

const (
	ColorMode1Bit ColorMode = iota
	ColorMode2Bit
)

func (c ColorMode) Palette() color.Palette {
	switch c {
	case ColorMode1Bit:
		return imaging.Palette1Bit()
	case ColorMode2Bit:
		return imaging.Palette2Bit()
	default:
		panic("invalid color mode")
	}
}
