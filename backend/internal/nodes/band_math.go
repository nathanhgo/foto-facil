package nodes

import (
	"errors"
	"image"
	"image/color"
)

// BandMathNode computes a generic per-pixel index between two grayscale
// "bands". With the default normalized-difference operation, (A-B)/(A+B),
// it reproduces the formula behind NDVI (Normalized Difference Vegetation
// Index) — the most common spectral index in remote sensing — once a real
// near-infrared band is available (see the GeoTIFF item in mvp.md). The
// index, which ranges from -1 to 1, is rescaled to 0-255 for visualization.
//
// ctx.Images[0] is band A, ctx.Images[1] is band B.
type BandMathNode struct {
	ID        string
	Operation string // "normalized_difference" (default), "difference" or "sum"
}

func (n *BandMathNode) GetID() string { return n.ID }

func (n *BandMathNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) < 2 {
		return errors.New("BandMathNode requires two input bands (A and B)")
	}

	op := n.Operation
	if op == "" {
		op = "normalized_difference"
	}

	bandA := toGrayValues(ctx.Images[0])
	bandB := toGrayValues(ctx.Images[1])

	w := imin(bandA.w, bandB.w)
	h := imin(bandA.h, bandB.h)

	dst := image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			av, _ := bandA.at(x, y)
			bv, _ := bandB.at(x, y)
			a := float64(av)
			b := float64(bv)

			var result float64
			switch op {
			case "difference":
				result = clampFloat((a - b + 255) / 2)
			case "sum":
				result = clampFloat((a + b) / 2)
			default: // normalized_difference
				var index float64
				if denom := a + b; denom > 0 {
					index = (a - b) / denom // ranges [-1, 1]
				}
				result = clampFloat((index + 1) / 2 * 255)
			}

			dst.Set(x, y, color.Gray{Y: uint8(result)})
		}
	}

	ctx.Images = []image.Image{dst}
	return nil
}
