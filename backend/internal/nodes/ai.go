package nodes

import (
	"errors"
	"image"
	"image/color"
	"math"
)

// AINode simulates smart ML tools like background removal and intelligence enhancement
type AINode struct {
	ID        string
	Tool      string // "background_removal" or "contrast_boost"
	Tolerance int    // Tolerance threshold for color distance (0-255)
}

func (n *AINode) GetID() string { return n.ID }

func (n *AINode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("AINode requires at least one input image")
	}

	tol := n.Tolerance
	if tol <= 0 {
		tol = 80 // Default to 80 for reasonable chroma key gradient handling
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		bounds := img.Bounds()
		
		if n.Tool == "background_removal" {
			// Sample 4 corners as references for background color
			c1R, c1G, c1B, _ := img.At(bounds.Min.X, bounds.Min.Y).RGBA()
			c2R, c2G, c2B, _ := img.At(bounds.Max.X-1, bounds.Min.Y).RGBA()
			c3R, c3G, c3B, _ := img.At(bounds.Min.X, bounds.Max.Y-1).RGBA()
			c4R, c4G, c4B, _ := img.At(bounds.Max.X-1, bounds.Max.Y-1).RGBA()

			corners := [4][3]uint8{
				{uint8(c1R >> 8), uint8(c1G >> 8), uint8(c1B >> 8)},
				{uint8(c2R >> 8), uint8(c2G >> 8), uint8(c2B >> 8)},
				{uint8(c3R >> 8), uint8(c3G >> 8), uint8(c3B >> 8)},
				{uint8(c4R >> 8), uint8(c4G >> 8), uint8(c4B >> 8)},
			}

			dst := image.NewRGBA(bounds)
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c := img.At(x, y)
					r, g, b, a := c.RGBA()
					r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

					// Distance between the current pixel and the closest corner color
					minDist := 999.0
					for _, corner := range corners {
						dist := math.Sqrt(math.Pow(float64(r8)-float64(corner[0]), 2) +
							math.Pow(float64(g8)-float64(corner[1]), 2) +
							math.Pow(float64(b8)-float64(corner[2]), 2))
						if dist < minDist {
							minDist = dist
						}
					}

					if minDist < float64(tol) {
						// Set transparent alpha
						dst.Set(x, y, color.RGBA{r8, g8, b8, 0})
					} else {
						dst.Set(x, y, color.RGBA{r8, g8, b8, a8})
					}
				}
			}
			processed = append(processed, dst)
		} else if n.Tool == "contrast_boost" {
			// Sigmoid contrast boost simulation
			dst := image.NewRGBA(bounds)
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c := img.At(x, y)
					r, g, b, a := c.RGBA()
					r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

					enhance := func(v uint8) uint8 {
						fv := float64(v) / 255.0
						res := 1.0 / (1.0 + math.Exp(-10.0*(fv-0.5)))
						return uint8(res * 255.0)
					}
					dst.Set(x, y, color.RGBA{enhance(r8), enhance(g8), enhance(b8), a8})
				}
			}
			processed = append(processed, dst)
		} else {
			// Fallback: pass-through
			processed = append(processed, img)
		}
	}

	ctx.Images = processed
	return nil
}
