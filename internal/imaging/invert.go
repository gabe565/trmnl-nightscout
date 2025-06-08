package imaging

import (
	"errors"
	"image"
)

var ErrNot1Bit = errors.New("image is not 1-bit")

func Invert1Bit(img *image.Paletted) error {
	if len(img.Palette) != 2 {
		return ErrNot1Bit
	}

	for i, c := range img.Pix {
		img.Pix[i] = 1 - c
	}
	return nil
}
