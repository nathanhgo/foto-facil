package nodes

import (
	"errors"
	"image"
	"image/color"
)

// ConvolutionNode applies a user-defined NxN kernel to every pixel — the
// fundamental operation behind spatial filters (sharpen, blur, edge
// detection) and every convolutional neural network layer. Border pixels
// are handled by replicating the edge (clamp-to-edge).
type ConvolutionNode struct {
	ID         string
	KernelSize int       // e.g. 3 for a 3x3 kernel
	Kernel     []float64 // Row-major, length KernelSize*KernelSize
}

func (n *ConvolutionNode) GetID() string { return n.ID }

func (n *ConvolutionNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("ConvolutionNode requires at least one input image")
	}

	size := n.KernelSize
	kernel := n.Kernel
	if size <= 0 || len(kernel) != size*size {
		size = 3
		kernel = []float64{0, 0, 0, 0, 1, 0, 0, 0, 0} // Identity kernel fallback
	}
	radius := size / 2

	var sum float64
	for _, k := range kernel {
		sum += k
	}
	if sum == 0 {
		sum = 1
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()
		dst := image.NewRGBA(image.Rect(0, 0, w, h))

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				var accR, accG, accB float64
				for ky := 0; ky < size; ky++ {
					for kx := 0; kx < size; kx++ {
						sx := clampInt(x+kx-radius, 0, w-1)
						sy := clampInt(y+ky-radius, 0, h-1)
						r, g, bch, _ := img.At(b.Min.X+sx, b.Min.Y+sy).RGBA()
						weight := kernel[ky*size+kx]
						accR += float64(r>>8) * weight
						accG += float64(g>>8) * weight
						accB += float64(bch>>8) * weight
					}
				}
				_, _, _, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
				dst.Set(x, y, color.RGBA{
					uint8(clampFloat(accR / sum)),
					uint8(clampFloat(accG / sum)),
					uint8(clampFloat(accB / sum)),
					uint8(a >> 8),
				})
			}
		}
		processed = append(processed, dst)
	}

	ctx.Images = processed
	return nil
}
