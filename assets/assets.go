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
	//go:embed Inter_18pt-Light.ttf
	InterLight []byte
	//go:embed Inter_18pt-Regular.ttf
	InterRegular []byte
	//go:embed Inter_18pt-SemiBold.ttf
	InterSemiBold []byte
	//go:embed OpenArrow-Regular.otf
	OpenArrow []byte
)

func Nightscout() image.Image {
	img := must.Must2(png.Decode(bytes.NewReader(nightscout)))
	return img
}
