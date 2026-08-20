package nodes

import (
	"errors"
	"image"
	"image/color"
)

// StackingNode averages N frames of the same scene to reduce random sensor
// noise (signal-to-noise ratio improves roughly with sqrt(N)). This is the
// classic astrophotography/satellite-calibration technique used to clean up
// noisy single exposures from star trackers and onboard cameras.
type StackingNode struct {
	ID string
}

func (n *StackingNode) GetID() string { return n.ID }

func (n *StackingNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) < 2 {
		return errors.New("StackingNode requires at least two input images to average")
	}

	base := ctx.Images[0].Bounds()
	w, h := base.Dx(), base.Dy()

	sumR := make([]float64, w*h)
	sumG := make([]float64, w*h)
	sumB := make([]float64, w*h)
	sumA := make([]float64, w*h)
	count := 0

	for _, img := range ctx.Images {
		b := img.Bounds()
		if b.Dx() != w || b.Dy() != h {
			continue // Skip frames with mismatched dimensions instead of failing the whole batch
		}
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, bch, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
				idx := y*w + x
				sumR[idx] += float64(r >> 8)
				sumG[idx] += float64(g >> 8)
				sumB[idx] += float64(bch >> 8)
				sumA[idx] += float64(a >> 8)
			}
		}
		count++
	}

	if count == 0 {
		return errors.New("StackingNode: no input frames matched the reference dimensions")
	}

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			dst.Set(x, y, color.RGBA{
				uint8(sumR[idx] / float64(count)),
				uint8(sumG[idx] / float64(count)),
				uint8(sumB[idx] / float64(count)),
				uint8(sumA[idx] / float64(count)),
			})
		}
	}

	ctx.Images = []image.Image{dst}
	return nil
}
