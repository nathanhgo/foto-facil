package nodes

import (
	"errors"
	"image"
	"image/color"
	"math"
)

// EdgeDetectionNode highlights intensity discontinuities using classic PDI
// operators: Sobel (gradient magnitude) or Laplacian (second derivative).
type EdgeDetectionNode struct {
	ID     string
	Method string // "sobel" (default) or "laplacian"
}

func (n *EdgeDetectionNode) GetID() string { return n.ID }

var sobelX = [3][3]float64{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
var sobelY = [3][3]float64{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}
var laplacianKernel = [3][3]float64{{0, 1, 0}, {1, -4, 1}, {0, 1, 0}}

func (n *EdgeDetectionNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("EdgeDetectionNode requires at least one input image")
	}

	method := n.Method
	if method == "" {
		method = "sobel"
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		gray := toGrayValues(img)
		at := func(x, y int) float64 {
			v, ok := gray.at(clampInt(x, 0, gray.w-1), clampInt(y, 0, gray.h-1))
			_ = ok
			return float64(v)
		}

		dst := image.NewGray(image.Rect(0, 0, gray.w, gray.h))
		for y := 0; y < gray.h; y++ {
			for x := 0; x < gray.w; x++ {
				var value float64
				if method == "laplacian" {
					var acc float64
					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {
							acc += at(x+kx, y+ky) * laplacianKernel[ky+1][kx+1]
						}
					}
					value = math.Abs(acc)
				} else {
					var gx, gy float64
					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {
							gx += at(x+kx, y+ky) * sobelX[ky+1][kx+1]
							gy += at(x+kx, y+ky) * sobelY[ky+1][kx+1]
						}
					}
					value = math.Sqrt(gx*gx + gy*gy)
				}
				dst.Set(x, y, color.Gray{Y: uint8(clampFloat(value))})
			}
		}
		processed = append(processed, dst)
	}

	ctx.Images = processed
	return nil
}
