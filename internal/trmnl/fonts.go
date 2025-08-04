package trmnl

import (
	"gabe565.com/trmnl-nightscout/assets"
	"gabe565.com/utils/must"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"gonum.org/v1/plot"
	plotfont "gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
)

//nolint:gochecknoglobals
var (
	light     *opentype.Font
	regular   *opentype.Font
	semiBold  *opentype.Font
	openArrow *opentype.Font
)

//nolint:gochecknoinits
func init() {
	light = must.Must2(opentype.Parse(assets.InterLight))
	regular = must.Must2(opentype.Parse(assets.InterRegular))
	semiBold = must.Must2(opentype.Parse(assets.InterSemiBold))
	openArrow = must.Must2(opentype.Parse(assets.OpenArrow))

	plotFont := plotfont.Font{Typeface: "Inter-SemiBold"}
	plotfont.DefaultCache.Add(plotfont.Collection{
		{Font: plotFont, Face: semiBold},
	})
	plot.DefaultFont = plotFont
	plotter.DefaultFont = plotFont
}

type Fonts struct {
	Reading font.Face
	Unit    font.Face
	Label   font.Face
	Info    font.Face
	Arrow   font.Face
}

func newFonts() Fonts {
	return Fonts{
		Reading: must.Must2(opentype.NewFace(light, &opentype.FaceOptions{
			Size: 74,
			DPI:  DPI,
		})),
		Unit: must.Must2(opentype.NewFace(light, &opentype.FaceOptions{
			Size: 23,
			DPI:  DPI,
		})),
		Label: must.Must2(opentype.NewFace(semiBold, &opentype.FaceOptions{
			Size: 11,
			DPI:  DPI,
		})),
		Info: must.Must2(opentype.NewFace(regular, &opentype.FaceOptions{
			Size: 23,
			DPI:  DPI,
		})),
		Arrow: must.Must2(opentype.NewFace(openArrow, &opentype.FaceOptions{
			Size: 20,
			DPI:  DPI,
		})),
	}
}
