package nodes

import (
	"image"
	"image/color"
	"testing"
)

// singleBrightPixelImage builds a mostly-dark image with one bright pixel at
// the center, useful for verifying erosion/dilation behavior.
func singleBrightPixelImage(size int) *image.RGBA {
	img := solidColorImage(size, size, 0, 0, 0)
	img.Set(size/2, size/2, color.RGBA{255, 255, 255, 255})
	return img
}

func TestMorphologyNodeErosionRemovesIsolatedBrightPixel(t *testing.T) {
	img := singleBrightPixelImage(5)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &MorphologyNode{ID: "morph1", Operation: "erosion", KernelSize: 3}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("MorphologyNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(2, 2).RGBA()
	if r>>8 != 0 {
		t.Errorf("expected erosion to remove the isolated bright pixel, got %d", r>>8)
	}
}

func TestMorphologyNodeDilationSpreadsBrightPixel(t *testing.T) {
	img := singleBrightPixelImage(5)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &MorphologyNode{ID: "morph2", Operation: "dilation", KernelSize: 3}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("MorphologyNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(1, 2).RGBA()
	if r>>8 != 255 {
		t.Errorf("expected dilation to spread brightness to neighboring pixel, got %d", r>>8)
	}
}

func TestMorphologyNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &MorphologyNode{ID: "morph3"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
