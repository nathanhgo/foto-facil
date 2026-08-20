package nodes

import (
	"errors"
	"image"
	"image/color"
)

// BandCompositeNode merges three grayscale "bands" into a single RGB image.
// This mirrors false-color composites used in remote sensing (e.g. stacking
// near-infrared, red and green bands to highlight vegetation), even though
// FotoFácil does not read real multispectral rasters yet (see stack.md for
// the GeoTIFF follow-up).
//
// ctx.Images[0], [1] and [2] are treated as the R, G and B bands respectively.
type BandCompositeNode struct {
	ID string
}

func (n *BandCompositeNode) GetID() string { return n.ID }

func (n *BandCompositeNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) < 3 {
		return errors.New("BandCompositeNode requires three input bands (R, G and B)")
	}

	bandR := toGrayValues(ctx.Images[0])
	bandG := toGrayValues(ctx.Images[1])
	bandB := toGrayValues(ctx.Images[2])

	w := imin(imin(bandR.w, bandG.w), bandB.w)
	h := imin(imin(bandR.h, bandG.h), bandB.h)

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, _ := bandR.at(x, y)
			g, _ := bandG.at(x, y)
			b, _ := bandB.at(x, y)
			dst.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	ctx.Images = []image.Image{dst}
	return nil
}
