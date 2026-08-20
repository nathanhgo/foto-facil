package nodes

import (
	"errors"
	"image"
	"image/color"
)

// ThresholdNode converts an image to pure black/white based on a luminance
// cutoff, either a fixed value ("global") or one computed automatically from
// the image's histogram via Otsu's method ("otsu").
type ThresholdNode struct {
	ID     string
	Method string // "global" (default) or "otsu"
	Value  int    // Cutoff (0-255) used when Method is "global"
}

func (n *ThresholdNode) GetID() string { return n.ID }

func (n *ThresholdNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("ThresholdNode requires at least one input image")
	}

	method := n.Method
	if method == "" {
		method = "global"
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		values := luminanceValues(img)
		cutoff := n.Value
		if method == "otsu" {
			cutoff = otsuThreshold(values)
		} else if cutoff <= 0 {
			cutoff = 128
		}

		b := img.Bounds()
		w, h := b.Dx(), b.Dy()
		dst := image.NewGray(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, bch, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
				lum := luminance8(r, g, bch)
				if int(lum) >= cutoff {
					dst.Set(x, y, color.Gray{Y: 255})
				} else {
					dst.Set(x, y, color.Gray{Y: 0})
				}
			}
		}
		processed = append(processed, dst)
	}

	ctx.Images = processed
	return nil
}

// otsuThreshold picks the cutoff that maximizes inter-class variance between
// the "below" and "above" pixel populations, per Otsu's 1979 method.
func otsuThreshold(values []float64) int {
	var hist [256]int
	for _, v := range values {
		hist[int(clampFloat(v))]++
	}
	total := len(values)
	if total == 0 {
		return 128
	}

	var sumAll float64
	for i, count := range hist {
		sumAll += float64(i * count)
	}

	var sumB, wB float64
	var best float64
	bestThreshold := 0

	for t := 0; t < 256; t++ {
		wB += float64(hist[t])
		if wB == 0 {
			continue
		}
		wF := float64(total) - wB
		if wF == 0 {
			break
		}
		sumB += float64(t * hist[t])
		mB := sumB / wB
		mF := (sumAll - sumB) / wF
		betweenVar := wB * wF * (mB - mF) * (mB - mF)
		if betweenVar > best {
			best = betweenVar
			bestThreshold = t
		}
	}

	return bestThreshold
}
