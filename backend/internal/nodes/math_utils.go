package nodes

import "image"

// imax returns the larger of two ints
func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// imin returns the smaller of two ints
func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// clampInt restricts v to the inclusive [lo, hi] range.
func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// clampFloat restricts v to the 0-255 byte range, ready to be cast to uint8.
func clampFloat(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

// luminance8 converts 16-bit RGBA channel samples (as returned by
// image.Color.RGBA()) into a single 8-bit luminance value using the
// standard perceptual weights (0.299R + 0.587G + 0.114B).
func luminance8(r, g, b uint32) uint8 {
	return uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
}

// grayBuffer is a flat 8-bit single-channel image buffer, used by nodes that
// treat inputs as independent "bands" (registration, band math/composite,
// morphology) instead of full RGBA images.
type grayBuffer struct {
	w, h int
	pix  []uint8
}

// toGrayValues converts an image to a luminance-only grayBuffer.
func toGrayValues(img image.Image) grayBuffer {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	pix := make([]uint8, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, bch, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			pix[y*w+x] = luminance8(r, g, bch)
		}
	}
	return grayBuffer{w: w, h: h, pix: pix}
}

// at safely reads a pixel, returning ok=false when (x, y) is out of bounds.
func (g grayBuffer) at(x, y int) (uint8, bool) {
	if x < 0 || y < 0 || x >= g.w || y >= g.h {
		return 0, false
	}
	return g.pix[y*g.w+x], true
}
