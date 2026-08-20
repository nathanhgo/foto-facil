package nodes

import (
	"image"
	"image/color"
	"testing"
)

func TestChangeDetectionNodeFlagsChangedRegion(t *testing.T) {
	before := solidColorImage(4, 2, 10, 10, 10)
	after := solidColorImage(4, 2, 10, 10, 10)
	// Change the right half to a bright color, simulating a cleared area
	// between two satellite passes over the same spot.
	for y := 0; y < 2; y++ {
		for x := 2; x < 4; x++ {
			after.Set(x, y, color.RGBA{240, 240, 240, 255})
		}
	}

	ctx := &ProcessContext{Images: []image.Image{before, after}}
	node := &ChangeDetectionNode{ID: "cd1", Threshold: 30}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("ChangeDetectionNode failed: %v", err)
	}
	if len(ctx.Images) != 1 {
		t.Fatalf("expected a single output image, got %d", len(ctx.Images))
	}

	result := ctx.Images[0]
	r, g, b, _ := result.At(3, 0).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Errorf("expected changed pixel to be highlighted in red, got (%d,%d,%d)", r>>8, g>>8, b>>8)
	}

	r, g, b, _ = result.At(0, 0).RGBA()
	if r>>8 == 255 && g>>8 == 0 && b>>8 == 0 {
		t.Errorf("expected unchanged pixel to NOT be highlighted in red")
	}
}

func TestChangeDetectionNodeRequiresTwoImages(t *testing.T) {
	ctx := &ProcessContext{Images: []image.Image{solidColorImage(2, 2, 1, 1, 1)}}
	node := &ChangeDetectionNode{ID: "cd2"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when fewer than two images are provided")
	}
}
