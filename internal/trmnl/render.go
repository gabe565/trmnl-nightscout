package trmnl

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"gabe565.com/trmnl-nightscout/assets"
	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/fetch"
	"gabe565.com/trmnl-nightscout/internal/imaging"
	"gabe565.com/utils/must"
	"git.sr.ht/~sbinet/gg"
	"github.com/makeworld-the-better-one/dither/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"gonum.org/v1/plot"
	plotfont "gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	vgdraw "gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

const (
	Width  = 800
	Height = 480
	DPI    = 124
	Margin = 25
)

//nolint:gochecknoglobals
var (
	light40    font.Face
	light128   font.Face
	regular39  font.Face
	semiBold20 font.Face
)

//nolint:gochecknoinits
func init() {
	plotFont := plotfont.Font{Typeface: "Inter-SemiBold"}
	plotfont.DefaultCache.Add(plotfont.Collection{
		{Font: plotFont, Face: must.Must2(opentype.Parse(assets.InterSemiBold))},
	})
	plot.DefaultFont = plotFont
	plotter.DefaultFont = plotFont

	light := must.Must2(opentype.Parse(assets.InterLight))
	light40 = must.Must2(opentype.NewFace(light, &opentype.FaceOptions{
		Size:    23,
		DPI:     DPI,
		Hinting: font.HintingFull,
	}))
	light128 = must.Must2(opentype.NewFace(light, &opentype.FaceOptions{
		Size:    74,
		DPI:     DPI,
		Hinting: font.HintingFull,
	}))

	regular := must.Must2(opentype.Parse(assets.InterRegular))
	regular39 = must.Must2(opentype.NewFace(regular, &opentype.FaceOptions{
		Size:    22.6,
		DPI:     DPI,
		Hinting: font.HintingFull,
	}))

	semiBold := must.Must2(opentype.Parse(assets.InterSemiBold))
	semiBold20 = must.Must2(opentype.NewFace(semiBold, &opentype.FaceOptions{
		Size:    11.5,
		DPI:     DPI,
		Hinting: font.HintingFull,
	}))
}

func Render(conf *config.Config, res *fetch.Response) (image.Image, error) {
	// Create regular image layer
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))

	// Create dithered layer
	dimg := image.NewRGBA(image.Rect(0, 0, Width, Height))
	draw.Draw(dimg, dimg.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	drawText(conf, res, img, dimg)
	drawPlot(conf, res, img, dimg)

	// Dither
	d := dither.NewDitherer([]color.Color{color.Black, color.White})
	d.Matrix = dither.FloydSteinberg
	d.Serpentine = true
	d.Dither(dimg)

	// Combine layers
	final := image.NewPaletted(img.Bounds(), color.Palette{color.Black, color.White})
	draw.Draw(final, final.Bounds(), dimg, image.Point{}, draw.Src)
	draw.Draw(final, final.Bounds(), img, image.Point{}, draw.Over)

	invert := conf.Invert
	if res.Properties.Bgnow.Last.Mgdl() <= conf.InvertBelow || res.Properties.Bgnow.Last.Mgdl() >= conf.InvertAbove {
		invert = !invert
	}
	if invert {
		if err := imaging.Invert1Bit(final); err != nil {
			return nil, err
		}
	}

	return final, nil
}

func drawText(conf *config.Config, res *fetch.Response, img, dimg *image.RGBA) {
	// Draw regular lines
	dc := gg.NewContextForRGBA(img)
	dc.SetDash(2, 4)
	dc.DrawLine(430, 113, 759, 113)
	dc.Stroke()
	dc.SetDash()

	drawer := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(color.Black),
	}

	// Last reading
	drawer.Face = light128
	const readingX, readingY = 49, 140
	drawer.Dot = fixed.P(readingX, readingY)
	drawer.DrawString(res.Properties.Bgnow.DisplayBg(conf.Units))

	if !res.Properties.IsRecent(conf.FetchDelay) {
		// Strikethrough
		dc.SetLineCapButt()
		dc.SetLineWidth(6)
		y := readingY - float64(light128.Metrics().XHeight)/64/2
		dc.DrawLine(readingX, y, float64(drawer.Dot.X)/64, y)
		dc.Stroke()
		dc.SetLineCapRound()
		dc.SetLineWidth(1)
	}

	drawer.Face = light40
	drawer.DrawString(" mg/dL")

	drawer.Face = semiBold20
	drawer.Dot = fixed.P(45, 170)
	drawer.DrawString("Last reading")

	// Updated
	drawer.Face = regular39
	drawer.Dot = fixed.P(450, 68)
	drawer.DrawString(res.Properties.Bgnow.Mills.Format("15:04"))

	drawer.Face = semiBold20
	drawer.Dot = fixed.P(450, 93)
	drawer.DrawString("Updated")

	// Nightscout logo
	nightscout := assets.Nightscout()
	draw.Draw(img, nightscout.Bounds().Add(image.Pt(630, 33)), nightscout, image.Point{}, draw.Over)

	// Direction
	drawer.Face = regular39
	drawer.Dot = fixed.P(450, 163)
	drawer.DrawString(res.Properties.Bgnow.Arrow())

	drawer.Face = semiBold20
	drawer.Dot = fixed.P(450, 183)
	drawer.DrawString("Direction")

	// Delta
	drawer.Face = regular39
	drawer.Dot = fixed.P(640, 163)
	drawer.DrawString(res.Properties.Delta.Display(conf.Units))

	drawer.Face = semiBold20
	drawer.Dot = fixed.P(640, 183)
	drawer.DrawString("Delta")

	// Draw dithered lines
	dc = gg.NewContextForRGBA(dimg)
	dc.SetColor(color.Gray{Y: 0xF2})
	dc.DrawRoundedRectangle(25, 30, 10, 150, 5)
	dc.Fill()
	dc.DrawRoundedRectangle(430, 30, 10, 70, 5)
	dc.Fill()
	dc.DrawRoundedRectangle(620, 30, 10, 70, 5)
	dc.Fill()
	dc.DrawRoundedRectangle(430, 125, 10, 70, 5)
	dc.Fill()
	dc.DrawRoundedRectangle(620, 125, 10, 70, 5)
	dc.Fill()
}

func drawPlot(conf *config.Config, res *fetch.Response, img, dimg *image.RGBA) {
	p := plot.New()
	p.BackgroundColor = color.Transparent

	p.Y.Min = 40
	p.Y.Max = 300
	p.Y.Padding = 0
	p.Y.Tick.Label.Font.Size = 10.8

	now := time.Now()
	p.X.Min = float64(now.Add(-conf.GraphDuration).Unix())
	p.X.Max = float64(now.Unix())
	p.X.Padding = 0
	p.X.Tick.Label.Font.Size = 10.8
	p.X.Tick.Marker = plot.TickerFunc(Ticks)

	// Render numbers and axes to non-dithered layer
	p.X.Color = color.Transparent
	p.X.Tick.Color = color.Transparent
	p.Y.Color = color.Transparent
	p.Y.Tick.Color = color.Transparent

	plotW := vg.Length(Width-2*Margin) * vg.Inch / DPI
	plotH := vg.Length(Height/2) * vg.Inch / DPI

	c := vgimg.NewWith(vgimg.UseWH(plotW, plotH), vgimg.UseDPI(DPI), vgimg.UseBackgroundColor(color.Transparent))
	p.Draw(vgdraw.New(c))
	draw.Draw(img, img.Bounds().Add(image.Pt(Margin, Height/2-Margin)), c.Image(), image.Point{}, draw.Over)

	p.X.Color = color.Black
	p.X.Tick.Color = color.Black
	p.Y.Color = color.Black
	p.Y.Tick.Color = color.Black

	p.Add(
		// Low threshold
		&plotter.Line{
			XYs: plotter.XYs{
				{X: p.X.Min, Y: conf.LowThreshold},
				{X: p.X.Max, Y: conf.LowThreshold},
			},
			LineStyle: vgdraw.LineStyle{
				Color:  color.Black,
				Width:  1,
				Dashes: []vg.Length{4, 2},
			},
			FillColor: color.Gray{Y: 0xF9},
		},

		// High threshold
		&plotter.Line{
			XYs: plotter.XYs{
				{X: p.X.Min, Y: conf.HighThreshold},
				{X: p.X.Max, Y: conf.HighThreshold},
			},
			LineStyle: vgdraw.LineStyle{
				Color:  color.Black,
				Width:  1,
				Dashes: []vg.Length{4, 2},
			},
		},

		&plotter.Polygon{
			XYs: []plotter.XYs{{
				{X: p.X.Min, Y: conf.HighThreshold},
				{X: p.X.Max, Y: conf.HighThreshold},
				{X: p.X.Max, Y: p.Y.Max},
				{X: p.X.Min, Y: p.Y.Max},
			}},
			Color: color.Gray{Y: 0xFB},
		},

		// Grid
		&plotter.Grid{
			Vertical: vgdraw.LineStyle{
				Color:  color.Black,
				Width:  1,
				Dashes: []vg.Length{1, 5},
			},
			Horizontal: vgdraw.LineStyle{
				Color:  color.Black,
				Width:  1,
				Dashes: []vg.Length{1, 5},
			},
		},
	)

	// Points
	points := make(plotter.XYs, 0, len(res.Entries))
	for _, entry := range res.Entries {
		points = append(points, plotter.XY{
			X: float64(entry.Date.Unix()),
			Y: float64(entry.SGV.Mgdl()),
		})
	}

	p.Add(&plotter.Scatter{
		XYs: points,
		GlyphStyle: vgdraw.GlyphStyle{
			Color:  color.Black,
			Radius: 2,
			Shape:  vgdraw.CircleGlyph{},
		},
	})

	p.X.Tick.Label.Color = color.Transparent
	p.Y.Tick.Label.Color = color.Transparent

	// Draw dithered plot parts
	c = vgimg.NewWith(vgimg.UseWH(plotW, plotH), vgimg.UseDPI(DPI))
	p.Draw(vgdraw.New(c))
	draw.Draw(dimg, dimg.Bounds().Add(image.Pt(Margin, Height/2-Margin)), c.Image(), image.Point{}, draw.Src)
}

func Ticks(min, max float64) []plot.Tick { //nolint:revive,predeclared
	start := time.Unix(int64(min), 0).Round(15 * time.Minute)
	end := time.Unix(int64(max), 0).Round(15 * time.Minute)

	ticks := make([]plot.Tick, 0, int(4*end.Sub(start).Hours()+1))
	for t := start; !t.After(end); t = t.Add(15 * time.Minute) {
		tick := plot.Tick{Value: float64(t.Unix())}
		if t.Minute() == 0 {
			tick.Label = t.Format("15:00")
		}
		ticks = append(ticks, tick)
	}
	return ticks
}
