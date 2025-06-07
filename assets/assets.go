package assets

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"

	"gabe565.com/utils/must"
)

//go:generate ./convert-icon.sh src/nightscout.svg dist/nightscout.png
var (
	//go:embed dist/nightscout.png
	nightscout []byte
)

func Nightscout() image.Image {
	img := must.Must2(png.Decode(bytes.NewReader(nightscout)))
	return img
}
