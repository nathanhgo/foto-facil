package nodes

import (
	"errors"
	"image"
	"image/color"
)

// HistogramEqualizationNode improves global contrast by redistributing pixel
// intensities so the cumulative histogram becomes approximately linear. Each
// color channel is equalized independently — a common simplification that
// can shift color balance slightly compared to luminance-only equalization,
// but is simple, fast and testable.
type HistogramEqualizationNode struct {
	ID string
}

func (n *HistogramEqualizationNode) GetID() string { return n.ID }

func (n *HistogramEqualizationNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("HistogramEqualizationNode requires at least one input image")
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()
		totalPixels := w * h

		var histR, histG, histB [256]int
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				r, g, bch, _ := img.At(x, y).RGBA()
				histR[r>>8]++
				histG[g>>8]++
				histB[bch>>8]++
			}
		}

		mapR := equalizationLUT(histR, totalPixels)
		mapG := equalizationLUT(histG, totalPixels)
		mapB := equalizationLUT(histB, totalPixels)

		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, bch, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
				dst.Set(x, y, color.RGBA{mapR[r>>8], mapG[g>>8], mapB[bch>>8], uint8(a >> 8)})
			}
		}
		processed = append(processed, dst)
	}

	ctx.Images = processed
	return nil
}

// equalizationLUT builds a 0-255 lookup table from a channel histogram using
// the standard cumulative-distribution-function equalization formula.
func equalizationLUT(hist [256]int, totalPixels int) [256]uint8 {
	var lut [256]uint8
	if totalPixels == 0 {
		return lut
	}
	var cumulative int
	for i, count := range hist {
		cumulative += count
		lut[i] = uint8(clampFloat(float64(cumulative) * 255.0 / float64(totalPixels)))
	}
	return lut
}
