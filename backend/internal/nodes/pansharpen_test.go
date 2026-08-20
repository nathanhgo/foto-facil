package nodes

import (
	"image"
	"testing"
)

func TestPanSharpenNodeOutputMatchesPanResolution(t *testing.T) {
	colorImg := solidColorImage(2, 2, 100, 50, 25) // Low-res color image
	panImg := solidColorImage(4, 4, 150, 150, 150) // High-res panchromatic image (brighter -> sharpens up)

	ctx := &ProcessContext{Images: []image.Image{colorImg, panImg}}
	node := &PanSharpenNode{ID: "pan1"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("PanSharpenNode failed: %v", err)
	}
	if len(ctx.Images) != 1 {
		t.Fatalf("expected a single output image, got %d", len(ctx.Images))
	}

	out := ctx.Images[0]
	b := out.Bounds()
	if b.Dx() != 4 || b.Dy() != 4 {
		t.Errorf("expected output to match the panchromatic resolution (4x4), got %dx%d", b.Dx(), b.Dy())
	}

	// Pan is brighter than the color image's luminance, so the ratio should
	// brighten the red channel above its original 100.
	r, _, _, _ := out.At(0, 0).RGBA()
	if r>>8 <= 100 {
		t.Errorf("expected pan-sharpened red channel to be brighter than 100, got %d", r>>8)
	}
}

func TestPanSharpenNodeRequiresTwoImages(t *testing.T) {
	ctx := &ProcessContext{Images: []image.Image{solidColorImage(2, 2, 1, 1, 1)}}
	node := &PanSharpenNode{ID: "pan2"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when fewer than two images are provided")
	}
}
