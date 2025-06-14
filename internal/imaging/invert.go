package imaging

import "image"

func InvertPaletted(img *image.Paletted) {
	colors := uint8(len(img.Palette)) - 1 //nolint:gosec
	for i, c := range img.Pix {
		img.Pix[i] = colors - c
	}
}
