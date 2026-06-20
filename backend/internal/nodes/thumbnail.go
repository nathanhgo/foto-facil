package nodes

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/png"
)

// ThumbnailNode generates a base64 PNG thumbnail for preview in the frontend
type ThumbnailNode struct {
	ID        string
	MaxWidth  int
	MaxHeight int
	// Filled after processing
	Thumbnails []string
}

func (n *ThumbnailNode) GetID() string {
	return n.ID
}

// resizeToFit scales down an image preserving aspect ratio, up to MaxWidth/MaxHeight
func resizeToFit(img image.Image, maxW, maxH int) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	if w == 0 || h == 0 {
		return img
	}

	scaleW := float64(maxW) / float64(w)
	scaleH := float64(maxH) / float64(h)
	scale := scaleW
	if scaleH < scale {
		scale = scaleH
	}
	if scale >= 1 {
		return img // No need to upscale
	}

	newW := int(float64(w) * scale)
	newH := int(float64(h) * scale)
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			srcX := int(float64(x) / scale)
			srcY := int(float64(y) / scale)
			dst.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}
	return dst
}

func (n *ThumbnailNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("ThumbnailNode requires at least one input image")
	}

	maxW := n.MaxWidth
	if maxW == 0 {
		maxW = 200
	}
	maxH := n.MaxHeight
	if maxH == 0 {
		maxH = 200
	}

	n.Thumbnails = nil
	for _, img := range ctx.Images {
		thumb := resizeToFit(img, maxW, maxH)
		var buf bytes.Buffer
		if err := png.Encode(&buf, thumb); err != nil {
			return err
		}
		encoded := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
		n.Thumbnails = append(n.Thumbnails, encoded)
	}

	return nil
}
