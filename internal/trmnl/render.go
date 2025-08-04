package trmnl

import (
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"sync"
	"time"

	"gabe565.com/trmnl-nightscout/assets"
	"gabe565.com/trmnl-nightscout/internal/bg"
	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/fetch"
	"gabe565.com/trmnl-nightscout/internal/imaging"
	"github.com/makeworld-the-better-one/dither/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"gonum.org/v1/plot"
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

func NewRenderer(conf config.Render, res *fetch.Response) *Renderer {
	return &Renderer{conf: conf, res: res}
}

type Renderer struct {
	conf config.Render
	res  *fetch.Response
	img  *image.Paletted
	mu   sync.Mutex
}

func (r *Renderer) Render() (image.Image, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.img != nil {
		return r.img, nil
	}

	r.img = image.NewPaletted(image.Rect(0, 0, Width, Height), r.conf.ColorMode.Palette())
	draw.Draw(r.img, r.img.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	if bgnow := r.res.Properties.Bgnow.Last; bgnow <= r.conf.InvertBelow || bgnow >= r.conf.InvertAbove {
		r.conf.Invert = !r.conf.Invert
	}

	r.drawText()
	r.drawPlot()

	if r.conf.Invert {
		imaging.InvertPaletted(r.img)
	}

	return r.img, nil
}

func (r *Renderer) drawText() {
	fonts := newFonts()

	drawer := &font.Drawer{
		Dst: r.img,
		Src: image.NewUniform(color.Black),
	}

	var dots image.Image
	if r.conf.ColorMode == config.ColorMode2Bit {
		c := imaging.Gray2
		if r.conf.Invert {
			c = imaging.Gray1
		}
		dots = imaging.NewDots(image.Pt(1, 1), false).SetForeground(c)
	} else {
		dots = imaging.NewDots(image.Pt(3, 1), true)
	}

	// Last reading
	r.drawClamp(image.Rect(25, 30, 35, 196), dots)

	drawer.Face = fonts.Reading
	const readingX, readingY = 49, 149
	drawer.Dot = fixed.P(readingX, readingY)
	drawer.DrawString(r.res.Properties.Bgnow.DisplayBg(r.conf.Unit))

	if time.Since(r.res.Properties.Bgnow.Mills.Time) > 15*time.Minute {
		// Strikethrough
		const thickness = 7
		y := readingY - int(float64(fonts.Reading.Metrics().XHeight)/64/2) - thickness/2
		rect := image.Rect(readingX, y, int(drawer.Dot.X/64), y+thickness)
		draw.Draw(r.img, rect, image.NewUniform(color.Black), image.Point{}, draw.Over)
	}

	drawer.Face = fonts.Unit
	drawer.DrawString(" " + r.conf.Unit.String())

	drawer.Face = fonts.Label
	drawer.Dot = fixed.P(45, 179)
	drawer.DrawString("Last reading")

	// Updated
	r.drawSegment(fonts.Label, fonts.Info, image.Pt(440, 30), dots, "Updated",
		r.res.Properties.Bgnow.Mills.Format(r.conf.TimeFormat),
	)

	// Nightscout logo
	r.drawClamp(image.Rect(640, 30, 650, 100), dots)

	nightscout := assets.Nightscout()
	draw.Draw(r.img, nightscout.Bounds().Add(image.Pt(650, 33)), nightscout, image.Point{}, draw.Over)

	// Horizontal separator
	var horizontalSrc image.Image
	if r.conf.ColorMode == config.ColorMode2Bit {
		horizontalSrc = image.NewUniform(imaging.Gray2)
	} else {
		horizontalSrc = imaging.NewDots(image.Pt(4, 0), false)
	}
	draw.Draw(r.img,
		image.Rect(440, 113, Width-Margin, 114),
		horizontalSrc, image.Point{}, draw.Src,
	)

	r.drawSegment(fonts.Label, fonts.Arrow, image.Pt(440, 125), dots, directionLabel, r.res.Properties.Bgnow.Arrow())
	r.drawSegment(fonts.Label, fonts.Info, image.Pt(640, 125), dots, "Delta",
		r.res.Properties.Delta.Display(r.conf.Unit),
	)
}

const directionLabel = "Direction"

func (r *Renderer) drawClamp(rect image.Rectangle, dots image.Image) {
	draw.Draw(r.img, rect, dots, image.Pt(0, 1), draw.Src)
	r.img.Set(rect.Min.X, rect.Min.Y+1, color.White)
	r.img.Set(rect.Min.X, rect.Max.Y-1, color.White)
}

func (r *Renderer) drawSegment(labelFace, infoFace font.Face, p image.Point, dots image.Image, label, value string) {
	r.drawClamp(image.Rect(p.X, p.Y, p.X+10, p.Y+70), dots)

	drawer := &font.Drawer{
		Dst: r.img,
		Src: image.NewUniform(color.Black),
	}

	drawer.Face = labelFace
	drawer.Dot = fixed.P(p.X+20, p.Y+61)
	drawer.DrawString(label)

	drawer.Face = infoFace
	drawer.Dot = fixed.P(p.X+20, p.Y+38)
	drawer.DrawString(value)
}

func (r *Renderer) drawPlot() {
	const (
		plotW = vg.Length(Width-2*Margin) * vg.Inch / DPI
		plotH = vg.Length(Height/2) * vg.Inch / DPI
	)

	p := plot.New()
	p.BackgroundColor = color.Transparent

	p.Y.Min = r.conf.GraphMin.Value(r.conf.Unit)
	p.Y.Max = r.conf.GraphMax.Value(r.conf.Unit)
	p.Y.Padding = 0
	p.Y.Tick.Label.Font.Size = 10
	if r.conf.Unit == bg.Mmol {
		ticks := make(plot.ConstantTicks, 0,
			int(r.conf.GraphMax.Value(r.conf.Unit))-int(r.conf.GraphMin.Value(r.conf.Unit))+1,
		)
		for i := int(r.conf.GraphMin.Value(r.conf.Unit)); i <= int(r.conf.GraphMax.Value(r.conf.Unit)); i++ {
			tick := plot.Tick{Value: float64(i)}
			if i%2 == 0 {
				tick.Label = strconv.Itoa(i)
			}
			ticks = append(ticks, tick)
		}
		p.Y.Tick.Marker = ticks
	}

	end := time.Now()
	start := end.Add(-r.conf.GraphDuration)
	p.X.Min = float64(start.Unix())
	p.X.Max = float64(end.Unix())
	p.X.Padding = 0
	p.X.Tick.Label.Font.Size = 10
	p.X.Tick.Marker = Ticks(r.conf)

	gridLine := vgdraw.LineStyle{
		Color: imaging.Gray2,
		Width: 1.2,
	}
	thresholdLine := vgdraw.LineStyle{
		Color: imaging.Gray1,
		Width: 1.2,
	}

	if r.conf.ColorMode == config.ColorMode1Bit {
		gridLine = vgdraw.LineStyle{
			Color:  color.Black,
			Width:  1,
			Dashes: []vg.Length{1, 4},
		}
		thresholdLine = vgdraw.LineStyle{
			Color: color.Black,
			Width: 1,
		}
	}

	grid := &plotter.Grid{
		Vertical:   gridLine,
		Horizontal: gridLine,
	}

	highLine := &plotter.Line{
		XYs: plotter.XYs{
			{X: p.X.Min, Y: r.conf.HighThreshold.Value(r.conf.Unit)},
			{X: p.X.Max, Y: r.conf.HighThreshold.Value(r.conf.Unit)},
		},
		LineStyle: thresholdLine,
	}

	highBg := &plotter.Polygon{
		XYs: []plotter.XYs{{
			{X: p.X.Min, Y: r.conf.HighThreshold.Value(r.conf.Unit)},
			{X: p.X.Max, Y: r.conf.HighThreshold.Value(r.conf.Unit)},
			{X: p.X.Max, Y: p.Y.Max},
			{X: p.X.Min, Y: p.Y.Max},
		}},
		Color: color.Black,
	}

	// Low threshold
	lowLine := &plotter.Line{
		XYs: plotter.XYs{
			{X: p.X.Min, Y: r.conf.LowThreshold.Value(r.conf.Unit)},
			{X: p.X.Max, Y: r.conf.LowThreshold.Value(r.conf.Unit)},
		},
		LineStyle: thresholdLine,
	}

	lowBg := &plotter.Polygon{
		XYs: []plotter.XYs{{
			{X: p.X.Min, Y: p.Y.Min},
			{X: p.X.Max, Y: p.Y.Min},
			{X: p.X.Max, Y: r.conf.LowThreshold.Value(r.conf.Unit)},
			{X: p.X.Min, Y: r.conf.LowThreshold.Value(r.conf.Unit)},
		}},
		Color: color.Black,
	}

	// Points
	pointsXY := make(plotter.XYs, 0, len(r.res.Entries))
	for _, entry := range r.res.Entries {
		if entry.Date.Before(start) || entry.Date.After(end) {
			continue
		}
		reading := entry.SGV.Value(r.conf.Unit)
		reading = min(reading, r.conf.GraphMax.Value(r.conf.Unit))
		reading = max(reading, r.conf.GraphMin.Value(r.conf.Unit))
		pointsXY = append(pointsXY, plotter.XY{
			X: float64(entry.Date.Unix()),
			Y: reading,
		})
	}

	// Points
	pointStroke := &plotter.Scatter{
		XYs: pointsXY,
		GlyphStyle: vgdraw.GlyphStyle{
			Color:  color.White,
			Radius: r.conf.PointStrokeRadius,
			Shape:  vgdraw.CircleGlyph{},
		},
	}
	points := &plotter.Scatter{
		XYs: pointsXY,
		GlyphStyle: vgdraw.GlyphStyle{
			Color:  color.Black,
			Radius: 2,
			Shape:  vgdraw.CircleGlyph{},
		},
	}

	// Render images based on color mode
	plotBounds := r.img.Bounds().Add(image.Pt(Margin, Height/2-Margin))
	if r.conf.ColorMode == config.ColorMode2Bit {
		// Hide elements for the bg image
		p.X.Color = color.Transparent
		p.Y.Color = color.Transparent
		p.X.Tick.Color = color.Transparent
		p.Y.Tick.Color = color.Transparent
		p.X.Tick.Label.Color = color.Transparent
		p.Y.Tick.Label.Color = color.Transparent

		// Render high bg mask
		p.Add(highBg)
		c := vgimg.NewWith(vgimg.UseWH(plotW, plotH), vgimg.UseDPI(DPI), vgimg.UseBackgroundColor(color.Transparent))
		p.Draw(vgdraw.New(c))
		highMask := c.Image()
		highBg.XYs = nil

		// Render high bg dots from mask
		dots := imaging.NewDots(image.Pt(3, 1), true).SetForeground(imaging.Gray1)
		draw.DrawMask(r.img, plotBounds, dots, image.Point{}, highMask, image.Point{}, draw.Over)

		// Render low bg mask
		p.Add(lowBg)
		c = vgimg.NewWith(vgimg.UseWH(plotW, plotH), vgimg.UseDPI(DPI), vgimg.UseBackgroundColor(color.Transparent))
		p.Draw(vgdraw.New(c))
		lowMask := c.Image()
		lowBg.XYs = nil

		// Render low bg dots from mask
		dots.SetGap(image.Pt(1, 0), true).SetForeground(imaging.Gray2)
		draw.DrawMask(r.img, plotBounds, dots, image.Point{}, lowMask, image.Point{}, draw.Over)

		// Show elements for the fg image
		p.X.Color = color.Black
		p.Y.Color = color.Black
		p.X.Tick.Color = color.Black
		p.Y.Tick.Color = color.Black
		p.X.Tick.Label.Color = color.Black
		p.Y.Tick.Label.Color = color.Black

		// Render fg image
		p.Add(grid, highLine, lowLine, pointStroke, points)
		c = vgimg.NewWith(vgimg.UseWH(plotW, plotH), vgimg.UseDPI(DPI), vgimg.UseBackgroundColor(color.Transparent))
		p.Draw(vgdraw.New(c))
		fgImg := c.Image()
		draw.Draw(r.img, plotBounds, fgImg, image.Point{}, draw.Over)
	} else {
		// Hide elements for the high/low mask
		p.X.Color = color.Transparent
		p.Y.Color = color.Transparent
		p.X.Tick.Color = color.Transparent
		p.Y.Tick.Color = color.Transparent
		p.X.Tick.Label.Color = color.Transparent
		p.Y.Tick.Label.Color = color.Transparent

		// Render high/low mask
		p.Add(highBg, lowBg)
		c := vgimg.NewWith(vgimg.UseWH(plotW, plotH), vgimg.UseDPI(DPI), vgimg.UseBackgroundColor(color.Transparent))
		p.Draw(vgdraw.New(c))
		bgMask := c.Image()
		highBg.XYs = nil
		lowBg.XYs = nil

		// Show labels for the fg image
		p.X.Tick.Label.Color = color.Black
		p.Y.Tick.Label.Color = color.Black

		// Render fg image
		c = vgimg.NewWith(vgimg.UseWH(plotW, plotH), vgimg.UseDPI(DPI), vgimg.UseBackgroundColor(color.Transparent))
		p.Draw(vgdraw.New(c))
		fgImg := c.Image()

		// Show/hide elements for dithered bg image
		p.X.Color = color.Black
		p.Y.Color = color.Black
		p.X.Tick.Color = color.Black
		p.Y.Tick.Color = color.Black
		p.X.Tick.Label.Color = color.Transparent
		p.Y.Tick.Label.Color = color.Transparent

		// Create the bg image
		c = vgimg.NewWith(vgimg.UseWH(plotW, plotH), vgimg.UseDPI(DPI))
		bgImg := c.Image()

		// Render the dots from mask
		dots := imaging.NewDots(image.Pt(3, 1), true)
		draw.DrawMask(bgImg, bgImg.Bounds(), dots, image.Point{}, bgMask, image.Point{}, draw.Over)

		// Render the plot
		p.Add(highBg, lowBg, grid, highLine, lowLine, pointStroke, points)
		p.Draw(vgdraw.New(c))

		// Dither the bg image
		d := dither.NewDitherer(r.conf.ColorMode.Palette())
		d.Matrix = dither.FloydSteinberg
		d.Serpentine = true
		d.Dither(bgImg)

		// Combine layers
		draw.Draw(r.img, plotBounds, bgImg, image.Point{}, draw.Src)
		draw.Draw(r.img, plotBounds, fgImg, image.Point{}, draw.Over)
	}
}

func Ticks(conf config.Render) plot.TickerFunc {
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
