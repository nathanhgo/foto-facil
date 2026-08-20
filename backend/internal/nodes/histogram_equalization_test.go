package nodes

import (
	"image"
	"image/color"
	"testing"
)

func TestEqualizationLUTStretchesBimodalHistogram(t *testing.T) {
	var hist [256]int
	hist[50] = 50
	hist[200] = 50

	lut := equalizationLUT(hist, 100)

	if lut[50] >= 130 {
		t.Errorf("expected the dark cluster to map below the midpoint, got %d", lut[50])
	}
	if lut[200] != 255 {
		t.Errorf("expected the brightest cluster to map to 255, got %d", lut[200])
	}
}

func TestHistogramEqualizationNodeChangesContrast(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{50, 50, 50, 255})
	img.Set(1, 0, color.RGBA{200, 200, 200, 255})

	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &HistogramEqualizationNode{ID: "eq1"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("HistogramEqualizationNode failed: %v", err)
	}

	r1, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	r2, _, _, _ := ctx.Images[0].At(1, 0).RGBA()
	if r1>>8 >= r2>>8 {
		t.Errorf("expected the darker pixel to remain darker after equalization: got %d vs %d", r1>>8, r2>>8)
	}
	if r2>>8 != 255 {
		t.Errorf("expected the brightest pixel to map to 255, got %d", r2>>8)
	}
}

func TestHistogramEqualizationNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &HistogramEqualizationNode{ID: "eq2"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
