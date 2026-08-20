package nodes

import (
	"image"
	"image/color"
	"testing"
)

func TestNormalizationNodeMinMaxStretchesRange(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{50, 50, 50, 255})
	img.Set(1, 0, color.RGBA{200, 200, 200, 255})

	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &NormalizationNode{ID: "norm1", Method: "minmax"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("NormalizationNode failed: %v", err)
	}

	r1, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	r2, _, _, _ := ctx.Images[0].At(1, 0).RGBA()
	if r1>>8 != 0 {
		t.Errorf("expected the darkest pixel to map to 0, got %d", r1>>8)
	}
	if r2>>8 != 255 {
		t.Errorf("expected the brightest pixel to map to 255, got %d", r2>>8)
	}
}

func TestNormalizationNodeZScoreCentersOnMean(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{100, 100, 100, 255})
	img.Set(1, 0, color.RGBA{100, 100, 100, 255})

	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &NormalizationNode{ID: "norm2", Method: "zscore"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("NormalizationNode failed: %v", err)
	}

	// With zero variance, z-score should fall back to mid-gray (128) instead of dividing by zero.
	r, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	if r>>8 != 128 {
		t.Errorf("expected zero-variance channel to fall back to 128, got %d", r>>8)
	}
}

func TestNormalizationNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &NormalizationNode{ID: "norm3"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
