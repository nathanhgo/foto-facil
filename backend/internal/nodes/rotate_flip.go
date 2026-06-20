package nodes

import (
	"errors"
	"image"
)

// FlipAxis represents horizontal or vertical flip
type FlipAxis int

const (
	FlipHorizontal FlipAxis = iota
	FlipVertical
	FlipBoth
)

// RotateAngle represents supported rotation angles
type RotateAngle int

const (
	Rotate90  RotateAngle = 90
	Rotate180 RotateAngle = 180
	Rotate270 RotateAngle = 270
)

// RotateFlipNode rotates and/or flips images
type RotateFlipNode struct {
	ID     string
	Angle  RotateAngle // 0 = no rotation
	Flip   FlipAxis
	DoFlip bool
	DoRotate bool
}

func (n *RotateFlipNode) GetID() string { return n.ID }

func (n *RotateFlipNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("RotateFlipNode requires at least one input image")
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		result := img

		if n.DoRotate {
			result = rotateImage(result, n.Angle)
		}
		if n.DoFlip {
			result = flipImage(result, n.Flip)
		}
		processed = append(processed, result)
	}

	ctx.Images = processed
	return nil
}

func rotateImage(img image.Image, angle RotateAngle) image.Image {
	b := img.Bounds()
	w, h := b.Max.X-b.Min.X, b.Max.Y-b.Min.Y

	switch angle {
	case Rotate90:
		dst := image.NewRGBA(image.Rect(0, 0, h, w))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(h-1-y, x, img.At(b.Min.X+x, b.Min.Y+y))
			}
		}
		return dst
	case Rotate180:
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(w-1-x, h-1-y, img.At(b.Min.X+x, b.Min.Y+y))
			}
		}
		return dst
	case Rotate270:
		dst := image.NewRGBA(image.Rect(0, 0, h, w))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(y, w-1-x, img.At(b.Min.X+x, b.Min.Y+y))
			}
		}
		return dst
	}
	return img
}

func flipImage(img image.Image, axis FlipAxis) image.Image {
	b := img.Bounds()
	w, h := b.Max.X-b.Min.X, b.Max.Y-b.Min.Y
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var srcX, srcY int
			switch axis {
			case FlipHorizontal:
				srcX, srcY = w-1-x, y
			case FlipVertical:
				srcX, srcY = x, h-1-y
			case FlipBoth:
				srcX, srcY = w-1-x, h-1-y
			}
			dst.Set(x, y, img.At(b.Min.X+srcX, b.Min.Y+srcY))
		}
	}
	return dst
}
