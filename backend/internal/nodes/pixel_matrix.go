package nodes

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// pixelMatrixCellSize is the rendered size (in canvas pixels) of each cell
// in the numeric grid.
const pixelMatrixCellSize = 32

// PixelMatrixNode renders a small region of the image as a numeric grid of
// luminance values (0-255), making the "image as a matrix of numbers"
// concept tangible — a staple of introductory digital image processing
// classes.
type PixelMatrixNode struct {
	ID      string
	RegionX int
	RegionY int
	// RegionSize is the width/height (in source pixels) of the square
	// sub-region to display; default 8.
	RegionSize int
}

func (n *PixelMatrixNode) GetID() string { return n.ID }

func (n *PixelMatrixNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("PixelMatrixNode requires at least one input image")
	}

	size := n.RegionSize
	if size <= 0 {
		size = 8
	}

	img := ctx.Images[0]
	b := img.Bounds()
	startX := clampInt(n.RegionX, b.Min.X, imax(b.Min.X, b.Max.X-size))
	startY := clampInt(n.RegionY, b.Min.Y, imax(b.Min.Y, b.Max.Y-size))

	canvas := image.NewRGBA(image.Rect(0, 0, size*pixelMatrixCellSize, size*pixelMatrixCellSize))

	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			x := startX + col
			y := startY + row
			var lum uint8
			if x < b.Max.X && y < b.Max.Y {
				r, g, bch, _ := img.At(x, y).RGBA()
				lum = luminance8(r, g, bch)
			}

			cellX := col * pixelMatrixCellSize
			cellY := row * pixelMatrixCellSize
			shade := color.RGBA{lum, lum, lum, 255}
			for py := 0; py < pixelMatrixCellSize; py++ {
				for px := 0; px < pixelMatrixCellSize; px++ {
					canvas.Set(cellX+px, cellY+py, shade)
				}
			}

			textColor := color.RGBA{255, 255, 255, 255}
			if lum >= 128 {
				textColor = color.RGBA{0, 0, 0, 255}
			}
			drawPixelValue(canvas, fmt.Sprintf("%d", lum), cellX+4, cellY+pixelMatrixCellSize/2+4, textColor)
		}
	}

	ctx.Images = []image.Image{canvas}
	return nil
}

// drawPixelValue renders s onto dst at (x, y) using a small built-in bitmap
// font, avoiding the need for external font assets.
func drawPixelValue(dst *image.RGBA, s string, x, y int, c color.Color) {
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
}
