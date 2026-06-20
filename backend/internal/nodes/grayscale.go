package nodes

import (
	"errors"
	"image"
	"image/color"
)

// GrayscaleNode converts images to grayscale
type GrayscaleNode struct {
	ID string
}

func (n *GrayscaleNode) GetID() string {
	return n.ID
}

func (n *GrayscaleNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("GrayscaleNode requires at least one input image")
	}

	var processed []image.Image

	for _, img := range ctx.Images {
		bounds := img.Bounds()
		newImg := image.NewGray(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				// Luminosity formula: 0.299R + 0.587G + 0.114B
				lum := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
				newImg.Set(x, y, color.Gray{Y: lum})
			}
		}
		processed = append(processed, newImg)
	}

	ctx.Images = processed
	return nil
}
