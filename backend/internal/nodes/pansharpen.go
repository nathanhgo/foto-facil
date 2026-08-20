package nodes

import (
	"errors"
	"image"
	"image/color"
)

// PanSharpenNode fuses a low-resolution color image with a higher-resolution
// grayscale ("panchromatic") image, producing a sharp color composite. This
// is the classic pan-sharpening technique used by optical Earth-observation
// satellites (Landsat, CBERS) that carry both multispectral and panchromatic
// sensors. Uses a Brovey-style ratio: each color channel is scaled by
// pan_luminance / color_luminance.
//
// ctx.Images[0] is the low-res color image, ctx.Images[1] is the
// high-res panchromatic image. Output resolution matches the pan image.
type PanSharpenNode struct {
	ID string
}

func (n *PanSharpenNode) GetID() string { return n.ID }

func (n *PanSharpenNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) < 2 {
		return errors.New("PanSharpenNode requires a color image and a panchromatic image")
	}

	colorImg := ctx.Images[0]
	panImg := ctx.Images[1]

	panBounds := panImg.Bounds()
	w, h := panBounds.Dx(), panBounds.Dy()
	colorBounds := colorImg.Bounds()
	cw, ch := colorBounds.Dx(), colorBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Nearest-neighbor upsample of the low-res color pixel to the pan grid
			cx := x * cw / w
			cy := y * ch / h
			cr, cg, cb, ca := colorImg.At(colorBounds.Min.X+cx, colorBounds.Min.Y+cy).RGBA()
			r8, g8, b8, a8 := uint8(cr>>8), uint8(cg>>8), uint8(cb>>8), uint8(ca>>8)

			pr, pg, pb, _ := panImg.At(panBounds.Min.X+x, panBounds.Min.Y+y).RGBA()
			panLum := luminance8(pr, pg, pb)
			colorLum := luminance8(cr, cg, cb)

			ratio := 1.0
			if colorLum > 0 {
				ratio = float64(panLum) / float64(colorLum)
			}

			dst.Set(x, y, color.RGBA{
				uint8(clampFloat(float64(r8) * ratio)),
				uint8(clampFloat(float64(g8) * ratio)),
				uint8(clampFloat(float64(b8) * ratio)),
				a8,
			})
		}
	}

	ctx.Images = []image.Image{dst}
	return nil
}
