package nodes

import (
	"errors"
	"image"
	"image/color"
)

// MorphologyNode applies structuring-element operations (erosion, dilation,
// opening, closing) commonly used to clean up shapes in binary/grayscale
// images before contour or shape analysis.
type MorphologyNode struct {
	ID         string
	Operation  string // "erosion" (default), "dilation", "opening" or "closing"
	KernelSize int    // Square structuring element size; default 3
}

func (n *MorphologyNode) GetID() string { return n.ID }

func (n *MorphologyNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("MorphologyNode requires at least one input image")
	}

	size := n.KernelSize
	if size <= 0 {
		size = 3
	}
	op := n.Operation
	if op == "" {
		op = "erosion"
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		gray := toGrayValues(img)
		var result grayBuffer

		switch op {
		case "dilation":
			result = dilate(gray, size)
		case "opening":
			result = dilate(erode(gray, size), size)
		case "closing":
			result = erode(dilate(gray, size), size)
		default:
			result = erode(gray, size)
		}

		dst := image.NewGray(image.Rect(0, 0, result.w, result.h))
		for y := 0; y < result.h; y++ {
			for x := 0; x < result.w; x++ {
				v, _ := result.at(x, y)
				dst.Set(x, y, color.Gray{Y: v})
			}
		}
		processed = append(processed, dst)
	}

	ctx.Images = processed
	return nil
}

func erode(g grayBuffer, size int) grayBuffer {
	return morphFilter(g, size, func(a, b uint8) uint8 {
		if a < b {
			return a
		}
		return b
	}, 255)
}

func dilate(g grayBuffer, size int) grayBuffer {
	return morphFilter(g, size, func(a, b uint8) uint8 {
		if a > b {
			return a
		}
		return b
	}, 0)
}

// morphFilter slides a size x size window over g, combining neighborhood
// values with combine (min for erosion, max for dilation). identity is the
// starting accumulator value (255 for min, 0 for max).
func morphFilter(g grayBuffer, size int, combine func(a, b uint8) uint8, identity uint8) grayBuffer {
	radius := size / 2
	out := grayBuffer{w: g.w, h: g.h, pix: make([]uint8, g.w*g.h)}
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			acc := identity
			for ky := -radius; ky <= radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					if v, ok := g.at(x+kx, y+ky); ok {
						acc = combine(acc, v)
					}
				}
			}
			out.pix[y*g.w+x] = acc
		}
	}
	return out
}
