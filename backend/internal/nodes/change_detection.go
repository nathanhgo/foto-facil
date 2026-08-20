package nodes

import (
	"errors"
	"image"
	"image/color"
)

// ChangeDetectionNode compares two images of the same scene captured at
// different times and highlights the regions that changed. This mirrors the
// technique used in satellite-based deforestation monitoring (e.g. INPE's
// PRODES/DETER programs), where two passes over the same area are diffed.
// It expects exactly two input images: ctx.Images[0] is "before" and
// ctx.Images[1] is "after".
type ChangeDetectionNode struct {
	ID        string
	Threshold int // Minimum grayscale difference (0-255) to flag a pixel as changed; default 30
}

func (n *ChangeDetectionNode) GetID() string { return n.ID }

func (n *ChangeDetectionNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) < 2 {
		return errors.New("ChangeDetectionNode requires two input images (before and after)")
	}

	threshold := n.Threshold
	if threshold <= 0 {
		threshold = 30
	}

	before := ctx.Images[0]
	after := ctx.Images[1]
	bBounds := before.Bounds()
	aBounds := after.Bounds()

	w := imin(bBounds.Dx(), aBounds.Dx())
	h := imin(bBounds.Dy(), aBounds.Dy())

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			br, bg, bb, _ := before.At(bBounds.Min.X+x, bBounds.Min.Y+y).RGBA()
			ar, ag, ab, _ := after.At(aBounds.Min.X+x, aBounds.Min.Y+y).RGBA()

			beforeLum := luminance8(br, bg, bb)
			afterLum := luminance8(ar, ag, ab)

			diff := int(afterLum) - int(beforeLum)
			if diff < 0 {
				diff = -diff
			}

			if diff >= threshold {
				dst.Set(x, y, color.RGBA{255, 0, 0, 255}) // Highlight the changed area in red
			} else {
				dst.Set(x, y, color.RGBA{afterLum, afterLum, afterLum, 255}) // Dim unchanged area to grayscale
			}
		}
	}

	ctx.Images = []image.Image{dst}
	return nil
}
