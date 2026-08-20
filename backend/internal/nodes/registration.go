package nodes

import (
	"errors"
	"image"
)

// RegistrationNode aligns a "moving" image to a "reference" image by
// searching for the pixel translation that minimizes the mean absolute
// grayscale difference between them. This is a simplified stand-in for the
// image co-registration step that precedes change detection or frame
// stacking on real satellite imagery, which rarely arrives perfectly
// aligned between passes.
//
// ctx.Images[0] is treated as the reference; ctx.Images[1] is the moving
// image. Output is [reference, alignedMoving].
type RegistrationNode struct {
	ID       string
	MaxShift int // Maximum search radius in pixels along X and Y; default 10
}

func (n *RegistrationNode) GetID() string { return n.ID }

func (n *RegistrationNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) < 2 {
		return errors.New("RegistrationNode requires a reference image and a moving image")
	}

	maxShift := n.MaxShift
	if maxShift <= 0 {
		maxShift = 10
	}

	reference := ctx.Images[0]
	moving := ctx.Images[1]

	refGray := toGrayValues(reference)
	movGray := toGrayValues(moving)

	bestDX, bestDY := 0, 0
	bestScore := -1.0

	for dy := -maxShift; dy <= maxShift; dy++ {
		for dx := -maxShift; dx <= maxShift; dx++ {
			score := sadOverlap(refGray, movGray, dx, dy)
			if score < 0 {
				continue
			}
			if bestScore < 0 || score < bestScore {
				bestScore = score
				bestDX, bestDY = dx, dy
			}
		}
	}

	aligned := shiftImage(moving, bestDX, bestDY, reference.Bounds())
	ctx.Images = []image.Image{reference, aligned}
	return nil
}

// sadOverlap computes the mean absolute difference between ref and mov when
// mov is conceptually shifted by (dx, dy), evaluated only over the region
// where both buffers overlap. Returns -1 when there is no overlap at all.
func sadOverlap(ref, mov grayBuffer, dx, dy int) float64 {
	var sum float64
	var count int
	for y := 0; y < ref.h; y++ {
		my := y - dy
		for x := 0; x < ref.w; x++ {
			mx := x - dx
			mv, ok := mov.at(mx, my)
			if !ok {
				continue
			}
			rv := ref.pix[y*ref.w+x]
			diff := int(rv) - int(mv)
			if diff < 0 {
				diff = -diff
			}
			sum += float64(diff)
			count++
		}
	}
	if count == 0 {
		return -1
	}
	return sum / float64(count)
}

// shiftImage translates img by (dx, dy), producing a new image sized to
// match targetBounds; pixels with no source data are left transparent black.
func shiftImage(img image.Image, dx, dy int, targetBounds image.Rectangle) image.Image {
	w, h := targetBounds.Dx(), targetBounds.Dy()
	b := img.Bounds()
	sw, sh := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sx := x - dx
			sy := y - dy
			if sx >= 0 && sy >= 0 && sx < sw && sy < sh {
				dst.Set(x, y, img.At(b.Min.X+sx, b.Min.Y+sy))
			}
		}
	}
	return dst
}
