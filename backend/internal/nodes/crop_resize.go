package nodes

import (
	"errors"
	"image"
)

// CropResizeNode crops and/or resizes the image
type CropResizeNode struct {
	ID     string
	Width  int
	Height int
	// Crop region (if zero, uses full image)
	CropX int
	CropY int
	CropW int
	CropH int
}

func (n *CropResizeNode) GetID() string { return n.ID }

func (n *CropResizeNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("CropResizeNode requires at least one input image")
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		result := img

		// Step 1: Crop (if crop region is set)
		if n.CropW > 0 && n.CropH > 0 {
			b := img.Bounds()
			cropX := imax(0, imin(n.CropX, b.Max.X))
			cropY := imax(0, imin(n.CropY, b.Max.Y))
			cropW := imax(1, imin(n.CropW, b.Max.X-cropX))
			cropH := imax(1, imin(n.CropH, b.Max.Y-cropY))
			cropped := image.NewRGBA(image.Rect(0, 0, cropW, cropH))
			for y := 0; y < cropH; y++ {
				for x := 0; x < cropW; x++ {
					cropped.Set(x, y, img.At(cropX+x, cropY+y))
				}
			}
			result = cropped
		}

		// Step 2: Resize (if target dimensions set)
		if n.Width > 0 && n.Height > 0 {
			src := result
			bounds := src.Bounds()
			srcW := bounds.Max.X - bounds.Min.X
			srcH := bounds.Max.Y - bounds.Min.Y
			dst := image.NewRGBA(image.Rect(0, 0, n.Width, n.Height))
			for y := 0; y < n.Height; y++ {
				for x := 0; x < n.Width; x++ {
					srcX := x * srcW / n.Width
					srcY := y * srcH / n.Height
					dst.Set(x, y, src.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
				}
			}
			result = dst
		}

		processed = append(processed, result)
	}

	ctx.Images = processed
	return nil
}
