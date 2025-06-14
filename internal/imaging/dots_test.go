package imaging

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDots(t *testing.T) {
	type args struct {
		gap     image.Point
		stagger bool
	}
	tests := []struct {
		name string
		args args
		want *DotImage
	}{
		{"3 5", args{image.Pt(3, 5), false}, &DotImage{image.Pt(4, 6), 0, Palette1Bit()}},
		{"3 5 stagger", args{image.Pt(3, 5), true}, &DotImage{image.Pt(4, 6), 2, Palette1Bit()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dots := NewDots(tt.args.gap, tt.args.stagger)
			assert.Equal(t, tt.want, dots)
		})
	}
}

func TestDotImage_NoStagger(t *testing.T) {
	d := NewDots(image.Pt(3, 3), false)

	cases := []struct {
		x, y int
		want color.Color
	}{
		{0, 0, color.Black}, // dot row, aligned
		{1, 0, color.White}, // dot row, off grid
		{4, 0, color.Black}, // next block
		{0, 1, color.White}, // non-dot row
		{0, 4, color.Black}, // dot at next block of rows
		{3, 4, color.White}, // same row, off grid
	}

	for _, tc := range cases {
		assert.Equal(t, tc.want, d.At(tc.x, tc.y))
	}
}

func TestDotImage_Stagger(t *testing.T) {
	d := NewDots(image.Pt(3, 3), true) // spaceX=4, spaceY=4, XOffset=2

	cases := []struct {
		x, y int
		want color.Color
	}{
		// first dot-row (y=0): same as no-stagger
		{0, 0, color.Black},
		{1, 0, color.White},

		// second dot-row (y=4), shifted by 2
		{0, 4, color.White}, // (0-2)%4 != 0
		{2, 4, color.Black}, // (2-2)%4 == 0
		{6, 4, color.Black}, // (6-2)%4 == 0
		{3, 4, color.White}, // (3-2)%4 != 0

		// non-dot rows still all white
		{2, 2, color.White},
		{5, 3, color.White},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.want, d.At(tc.x, tc.y))
	}
}
