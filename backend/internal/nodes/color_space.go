package nodes

import (
	"errors"
	"image"
	"image/color"
	"math"
)

// ColorSpaceNode converts images between common color spaces used in image
// processing. Since the Go standard image types only render true RGB(A),
// non-RGB modes (HSV, Lab, YCbCr) are visualized by mapping their three
// components onto the R, G and B channels of the output image — a common
// technique to inspect each channel independently.
type ColorSpaceNode struct {
	ID   string
	Mode string // "grayscale" (default), "hsv", "lab" or "ycbcr"
}

func (n *ColorSpaceNode) GetID() string { return n.ID }

func (n *ColorSpaceNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("ColorSpaceNode requires at least one input image")
	}

	mode := n.Mode
	if mode == "" {
		mode = "grayscale"
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()

		if mode == "grayscale" {
			dst := image.NewGray(image.Rect(0, 0, w, h))
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r, g, bch, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
					dst.Set(x, y, color.Gray{Y: luminance8(r, g, bch)})
				}
			}
			processed = append(processed, dst)
			continue
		}

		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, bch, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
				r8, g8, b8 := float64(r>>8), float64(g>>8), float64(bch>>8)

				var c1, c2, c3 uint8
				switch mode {
				case "hsv":
					c1, c2, c3 = rgbToHSVChannels(r8, g8, b8)
				case "lab":
					c1, c2, c3 = rgbToLabChannels(r8, g8, b8)
				default: // ycbcr
					c1, c2, c3 = rgbToYCbCrChannels(r8, g8, b8)
				}

				dst.Set(x, y, color.RGBA{c1, c2, c3, uint8(a >> 8)})
			}
		}
		processed = append(processed, dst)
	}

	ctx.Images = processed
	return nil
}

func rgbToHSVChannels(r, g, b float64) (uint8, uint8, uint8) {
	r, g, b = r/255, g/255, b/255
	maxV := math.Max(r, math.Max(g, b))
	minV := math.Min(r, math.Min(g, b))
	delta := maxV - minV

	var h float64
	switch {
	case delta == 0:
		h = 0
	case maxV == r:
		h = 60 * math.Mod((g-b)/delta, 6)
	case maxV == g:
		h = 60 * ((b-r)/delta + 2)
	default:
		h = 60 * ((r-g)/delta + 4)
	}
	if h < 0 {
		h += 360
	}

	var s float64
	if maxV > 0 {
		s = delta / maxV
	}

	return uint8(clampFloat(h / 360 * 255)), uint8(clampFloat(s * 255)), uint8(clampFloat(maxV * 255))
}

func rgbToYCbCrChannels(r, g, b float64) (uint8, uint8, uint8) {
	y := 0.299*r + 0.587*g + 0.114*b
	cb := 128 - 0.168736*r - 0.331264*g + 0.5*b
	cr := 128 + 0.5*r - 0.418688*g - 0.081312*b
	return uint8(clampFloat(y)), uint8(clampFloat(cb)), uint8(clampFloat(cr))
}

func rgbToLabChannels(r, g, b float64) (uint8, uint8, uint8) {
	toLinear := func(v float64) float64 {
		v /= 255
		if v <= 0.04045 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	lr, lg, lb := toLinear(r), toLinear(g), toLinear(b)

	// Linear RGB -> XYZ (D65 illuminant)
	x := lr*0.4124 + lg*0.3576 + lb*0.1805
	y := lr*0.2126 + lg*0.7152 + lb*0.0722
	z := lr*0.0193 + lg*0.1192 + lb*0.9505

	x /= 0.95047
	z /= 1.08883

	f := func(t float64) float64 {
		if t > 0.008856 {
			return math.Cbrt(t)
		}
		return 7.787*t + 16.0/116.0
	}
	fx, fy, fz := f(x), f(y), f(z)

	l := 116*fy - 16
	aStar := 500 * (fx - fy)
	bStar := 200 * (fy - fz)

	// Scale L (0-100) and a*/b* (roughly -128..127) into 0-255 for visualization
	return uint8(clampFloat(l * 2.55)), uint8(clampFloat(aStar + 128)), uint8(clampFloat(bStar + 128))
}
