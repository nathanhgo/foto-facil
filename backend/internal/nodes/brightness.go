package nodes

import (
	"errors"
	"image"
	"image/color"
)

// BrightnessNode adjusts the brightness of images in the batch
type BrightnessNode struct {
	ID         string
	Brightness int // Value between -255 and 255
}

func (n *BrightnessNode) GetID() string {
	return n.ID
}

func (n *BrightnessNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("BrightnessNode requires at least one input image")
	}

	var processed []image.Image

	for _, img := range ctx.Images {
		bounds := img.Bounds()
		newImg := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()

				// Convert from 16-bit to 8-bit
				r8 := int(r >> 8)
				g8 := int(g >> 8)
				b8 := int(b >> 8)

				// Apply brightness
				r8 = clamp(r8 + n.Brightness)
				g8 = clamp(g8 + n.Brightness)
				b8 = clamp(b8 + n.Brightness)

				newImg.Set(x, y, color.RGBA{uint8(r8), uint8(g8), uint8(b8), uint8(a >> 8)})
			}
		}
		processed = append(processed, newImg)
	}

	ctx.Images = processed
	return nil
}

func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}
