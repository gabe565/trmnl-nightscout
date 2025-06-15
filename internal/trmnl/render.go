package trmnl

import (
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"time"

	"gabe565.com/trmnl-nightscout/assets"
	"gabe565.com/trmnl-nightscout/internal/bg"
	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/fetch"
	"gabe565.com/trmnl-nightscout/internal/imaging"
	"gabe565.com/utils/must"
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
	light23     font.Face
	light74     font.Face
	regular23   font.Face
	semiBold11  font.Face
	openArrow23 font.Face
)

//nolint:gochecknoinits
func init() {
	light := must.Must2(opentype.Parse(assets.InterLight))
	light23 = must.Must2(opentype.NewFace(light, &opentype.FaceOptions{
		Size:    23,
		DPI:     DPI,
		Hinting: font.HintingFull,
	}))
	light74 = must.Must2(opentype.NewFace(light, &opentype.FaceOptions{
		Size:    74,
		DPI:     DPI,
		Hinting: font.HintingFull,
	}))

	regular := must.Must2(opentype.Parse(assets.InterRegular))
	regular23 = must.Must2(opentype.NewFace(regular, &opentype.FaceOptions{
		Size:    23,
		DPI:     DPI,
		Hinting: font.HintingFull,
	}))

	semiBold := must.Must2(opentype.Parse(assets.InterSemiBold))
	semiBold11 = must.Must2(opentype.NewFace(semiBold, &opentype.FaceOptions{
		Size:    11,
		DPI:     DPI,
		Hinting: font.HintingFull,
	}))

	plotFont := plotfont.Font{Typeface: "Inter-SemiBold"}
	plotfont.DefaultCache.Add(plotfont.Collection{
		{Font: plotFont, Face: semiBold},
	})
	plot.DefaultFont = plotFont
	plotter.DefaultFont = plotFont

	openArrow := must.Must2(opentype.Parse(assets.OpenArrow))
	openArrow23 = must.Must2(opentype.NewFace(openArrow, &opentype.FaceOptions{
		Size:    20,
		DPI:     DPI,
		Hinting: font.HintingFull,
	}))
}

func Render(conf *config.Config, res *fetch.Response) (image.Image, error) {
	// Create regular image layer
	img := image.NewPaletted(image.Rect(0, 0, Width, Height), imaging.Palette1Bit())
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	drawText(conf, res, img)
	drawPlot(conf, res, img)

	invert := conf.Invert
	bgnow := res.Properties.Bgnow.Last.Value(conf.Units)
	if bgnow <= conf.InvertBelow || bgnow >= conf.InvertAbove {
		invert = !invert
	}
	if invert {
		imaging.InvertPaletted(img)
	}

	return img, nil
}

func drawText(conf *config.Config, res *fetch.Response, img *image.Paletted) {
	drawer := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(color.Black),
	}

	dots := imaging.NewDots(image.Pt(3, 1), true)

	// Last reading
	draw.Draw(img, image.Rect(25, 30, 35, 180), dots, image.Pt(0, 1), draw.Src)

	drawer.Face = light74
	const readingX, readingY = 49, 140
	drawer.Dot = fixed.P(readingX, readingY)
	drawer.DrawString(res.Properties.Bgnow.DisplayBg(conf.Units))

	if time.Since(res.Properties.Bgnow.Mills.Time) > 15*time.Minute {
		// Strikethrough
		const thickness = 7
		y := readingY - int(float64(light74.Metrics().XHeight)/64/2) - thickness/2
		rect := image.Rect(readingX, y, int(drawer.Dot.X/64), y+thickness)
		draw.Draw(img, rect, image.NewUniform(color.Black), image.Point{}, draw.Over)
	}

	drawer.Face = light23
	drawer.DrawString(" " + conf.Units.String())

	drawer.Face = semiBold11
	drawer.Dot = fixed.P(45, 170)
	drawer.DrawString("Last reading")

	// Updated
	drawSegment(img, image.Pt(440, 30), "Updated", res.Properties.Bgnow.Mills.Format(conf.TimeFormat))

	// Nightscout logo
	draw.Draw(img, image.Rect(640, 30, 650, 100), dots, image.Pt(0, 1), draw.Src)

	nightscout := assets.Nightscout()
	draw.Draw(img, nightscout.Bounds().Add(image.Pt(650, 33)), nightscout, image.Point{}, draw.Over)

	// Horizontal separator
	draw.Draw(img,
		image.Rect(440, 113, Width-Margin, 114),
		imaging.NewDots(image.Pt(4, 0), false),
		image.Point{}, draw.Src,
	)

	drawSegment(img, image.Pt(440, 125), directionLabel, res.Properties.Bgnow.Arrow())
	drawSegment(img, image.Pt(640, 125), "Delta", res.Properties.Delta.Display(conf.Units))
}

const directionLabel = "Direction"

func drawSegment(img *image.Paletted, p image.Point, label, value string) {
	dots := imaging.NewDots(image.Pt(3, 1), true)
	draw.Draw(img, image.Rect(p.X, p.Y, p.X+10, p.Y+70), dots, image.Pt(0, 1), draw.Src)

	drawer := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(color.Black),
	}

	drawer.Face = semiBold11
	drawer.Dot = fixed.P(p.X+20, p.Y+61)
	drawer.DrawString(label)

	drawer.Face = regular23
	if label == directionLabel {
		drawer.Face = openArrow23
	}
	drawer.Dot = fixed.P(p.X+20, p.Y+38)
	drawer.DrawString(value)
}

func drawPlot(conf *config.Config, res *fetch.Response, img *image.Paletted) {
	p := plot.New()
	p.BackgroundColor = color.Transparent

	p.Y.Min = float64(conf.GraphMin)
	p.Y.Max = float64(conf.GraphMax)
	p.Y.Padding = 0
	p.Y.Tick.Label.Font.Size = 10
	if conf.Units == bg.Mmol {
		ticks := make(plot.ConstantTicks, 0, conf.GraphMax-conf.GraphMin+1)
		for i := conf.GraphMin; i <= conf.GraphMax; i++ {
			tick := plot.Tick{Value: float64(i)}
			if i%2 == 0 {
				tick.Label = strconv.Itoa(i)
			}
			ticks = append(ticks, tick)
		}
		p.Y.Tick.Marker = ticks
	}

	now := time.Now()
	start := now.Add(-conf.GraphDuration)
	p.X.Min = float64(start.Unix())
	p.X.Max = float64(now.Unix())
	p.X.Padding = 0
	p.X.Tick.Label.Font.Size = 10
	p.X.Tick.Marker = Ticks(conf)

	// Render numbers and axes to non-dithered layer
	p.X.Color = color.Transparent
	p.X.Tick.Color = color.Transparent
	p.Y.Color = color.Transparent
	p.Y.Tick.Color = color.Transparent

	plotW := vg.Length(Width-2*Margin) * vg.Inch / DPI
	plotH := vg.Length(Height/2) * vg.Inch / DPI

	c := vgimg.NewWith(vgimg.UseWH(plotW, plotH), vgimg.UseDPI(DPI), vgimg.UseBackgroundColor(color.Transparent))
	p.Draw(vgdraw.New(c))
	axisImg := c.Image()

	p.X.Color = color.Black
	p.X.Tick.Color = color.Black
	p.Y.Color = color.Black
	p.Y.Tick.Color = color.Black

	p.Add(
		// High threshold background
		&plotter.Polygon{
			XYs: []plotter.XYs{{
				{X: p.X.Min, Y: conf.HighThreshold},
				{X: p.X.Max, Y: conf.HighThreshold},
				{X: p.X.Max, Y: p.Y.Max},
				{X: p.X.Min, Y: p.Y.Max},
			}},
			Color: color.Gray{Y: conf.HighBackgroundShade},
		},

		// Low threshold background
		&plotter.Polygon{
			XYs: []plotter.XYs{{
				{X: p.X.Min, Y: p.Y.Min},
				{X: p.X.Max, Y: p.Y.Min},
				{X: p.X.Max, Y: conf.LowThreshold},
				{X: p.X.Min, Y: conf.LowThreshold},
			}},
			Color: color.Gray{Y: conf.LowBackgroundShade},
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
		},
	)

	// Points
	points := make(plotter.XYs, 0, len(res.Entries))
	for _, entry := range res.Entries {
		if entry.Date.Before(start) {
			continue
		}
		reading := max(float64(conf.GraphMin), min(float64(conf.GraphMax), entry.SGV.Value(conf.Units)))
		points = append(points, plotter.XY{
			X: float64(entry.Date.Unix()),
			Y: reading,
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
	plotImg := c.Image()

	// Dither
	d := dither.NewDitherer(imaging.Palette1Bit())
	d.Matrix = dither.FloydSteinberg
	d.Serpentine = true
	d.Dither(plotImg)

	// Combine layers
	plotBounds := img.Bounds().Add(image.Pt(Margin, Height/2-Margin))
	draw.Draw(img, plotBounds, plotImg, image.Point{}, draw.Src)
	draw.Draw(img, plotBounds, axisImg, image.Point{}, draw.Over)
}

func Ticks(conf *config.Config) plot.TickerFunc {
	interval := 15 * time.Minute
	if conf.GraphDuration > 8*time.Hour {
		interval = 30 * time.Minute
	}

	var showEvery int
	switch {
	case conf.GraphDuration >= 18*time.Hour:
		showEvery = 3
	case conf.GraphDuration >= 10*time.Hour:
		showEvery = 2
	default:
		showEvery = 1
	}

	var lastMin, lastMax float64
	var ticks []plot.Tick

	return func(minVal, maxVal float64) []plot.Tick {
		if minVal == lastMin && maxVal == lastMax {
			return ticks
		}

		start := time.Unix(int64(minVal), 0).Round(interval)
		end := time.Unix(int64(maxVal), 0).Round(interval)

		ticks = make([]plot.Tick, 0, int(float64(time.Hour/interval)*end.Sub(start).Hours()+1))
		var hourIdx int

		for t := start; !t.After(end); t = t.Add(interval) {
			tick := plot.Tick{Value: float64(t.Unix())}

			if t.Minute() == 0 {
				if hourIdx%showEvery == 0 {
					tick.Label = t.Format(conf.TimeFormat)
				} else {
					tick.Label = " "
				}
				hourIdx++
			}

			ticks = append(ticks, tick)
		}

		lastMin, lastMax = minVal, maxVal
		return ticks
	}
}
