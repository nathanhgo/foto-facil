package nodes

import (
	"errors"
	"image"
	"image/color"
)

// CompareNode holds the original and processed images side-by-side for before/after preview
type CompareNode struct {
	ID string
}

func (n *CompareNode) GetID() string { return n.ID }

// Process expects an original image in ctx.OriginalImages and a processed image in ctx.Images
// It produces a side-by-side combined image so the thumbnail shows both halves
func (n *CompareNode) Process(ctx *ProcessContext) error {
	if len(ctx.OriginalImages) == 0 {
		return errors.New("CompareNode: no original images found in context")
	}
	if len(ctx.Images) == 0 {
		return errors.New("CompareNode: no processed images found in context")
	}

	left := ctx.OriginalImages[0]
	right := ctx.Images[len(ctx.Images)-1]

	lb := left.Bounds()
	rb := right.Bounds()

	// Use the smaller height
	h := lb.Max.Y - lb.Min.Y
	if rh := rb.Max.Y - rb.Min.Y; rh < h {
		h = rh
	}
	lw := lb.Max.X - lb.Min.X
	rw := rb.Max.X - rb.Min.X

	combined := image.NewRGBA(image.Rect(0, 0, lw+rw+2, h))

	// Draw left half
	for y := 0; y < h; y++ {
		for x := 0; x < lw; x++ {
			combined.Set(x, y, left.At(lb.Min.X+x, lb.Min.Y+y))
		}
		// Divider line (1px white)
		combined.Set(lw, y, color.RGBA{200, 200, 200, 255})
		combined.Set(lw+1, y, color.RGBA{200, 200, 200, 255})
		// Draw right half
		for x := 0; x < rw; x++ {
			combined.Set(lw+2+x, y, right.At(rb.Min.X+x, rb.Min.Y+y))
		}
	}

	// Replace the images buffer with the combined image so ThumbnailNode generates it
	ctx.Images = []image.Image{combined}
	return nil
}
