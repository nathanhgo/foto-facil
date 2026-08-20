package nodes

import (
	"errors"
	"image"
	"image/color"
	"math"
	"sort"
)

// BlurNode smooths images to reduce detail/noise using one of three classic
// spatial filters: Gaussian, mean (box) or median.
type BlurNode struct {
	ID         string
	Method     string // "gaussian" (default), "mean" or "median"
	KernelSize int    // Odd window size; default 3
}

func (n *BlurNode) GetID() string { return n.ID }

func (n *BlurNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("BlurNode requires at least one input image")
	}

	size := n.KernelSize
	if size <= 0 {
		size = 3
	}
	if size%2 == 0 {
		size++ // Kernels must be odd-sized to have a well-defined center
	}
	method := n.Method
	if method == "" {
		method = "gaussian"
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()
		dst := image.NewRGBA(image.Rect(0, 0, w, h))

		switch method {
		case "median":
			applyMedianBlur(img, dst, w, h, size)
		case "mean":
			applyWeightedBlur(img, dst, w, h, boxKernel(size))
		default:
			applyWeightedBlur(img, dst, w, h, gaussianKernel(size))
		}

		processed = append(processed, dst)
	}

	ctx.Images = processed
	return nil
}

func boxKernel(size int) [][]float64 {
	k := make([][]float64, size)
	w := 1.0 / float64(size*size)
	for i := range k {
		k[i] = make([]float64, size)
		for j := range k[i] {
			k[i][j] = w
		}
	}
	return k
}

func gaussianKernel(size int) [][]float64 {
	sigma := float64(size) / 3.0
	if sigma <= 0 {
		sigma = 1
	}
	radius := size / 2
	k := make([][]float64, size)
	var sum float64
	for y := -radius; y <= radius; y++ {
		row := make([]float64, size)
		for x := -radius; x <= radius; x++ {
			v := math.Exp(-float64(x*x+y*y) / (2 * sigma * sigma))
			row[x+radius] = v
			sum += v
		}
		k[y+radius] = row
	}
	for y := range k {
		for x := range k[y] {
			k[y][x] /= sum
		}
	}
	return k
}

func applyWeightedBlur(img image.Image, dst *image.RGBA, w, h int, kernel [][]float64) {
	b := img.Bounds()
	size := len(kernel)
	radius := size / 2
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var accR, accG, accB float64
			for ky := 0; ky < size; ky++ {
				for kx := 0; kx < size; kx++ {
					sx := clampInt(x+kx-radius, 0, w-1)
					sy := clampInt(y+ky-radius, 0, h-1)
					r, g, bch, _ := img.At(b.Min.X+sx, b.Min.Y+sy).RGBA()
					weight := kernel[ky][kx]
					accR += float64(r>>8) * weight
					accG += float64(g>>8) * weight
					accB += float64(bch>>8) * weight
				}
			}
			_, _, _, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			dst.Set(x, y, color.RGBA{uint8(clampFloat(accR)), uint8(clampFloat(accG)), uint8(clampFloat(accB)), uint8(a >> 8)})
		}
	}
}

func applyMedianBlur(img image.Image, dst *image.RGBA, w, h, size int) {
	b := img.Bounds()
	radius := size / 2
	windowLen := size * size
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			rs := make([]int, 0, windowLen)
			gs := make([]int, 0, windowLen)
			bs := make([]int, 0, windowLen)
			for ky := -radius; ky <= radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					sx := clampInt(x+kx, 0, w-1)
					sy := clampInt(y+ky, 0, h-1)
					r, g, bch, _ := img.At(b.Min.X+sx, b.Min.Y+sy).RGBA()
					rs = append(rs, int(r>>8))
					gs = append(gs, int(g>>8))
					bs = append(bs, int(bch>>8))
				}
			}
			sort.Ints(rs)
			sort.Ints(gs)
			sort.Ints(bs)
			mid := windowLen / 2
			_, _, _, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			dst.Set(x, y, color.RGBA{uint8(rs[mid]), uint8(gs[mid]), uint8(bs[mid]), uint8(a >> 8)})
		}
	}
}
