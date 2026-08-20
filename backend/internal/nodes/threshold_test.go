package nodes

import (
	"image"
	"image/color"
	"testing"
)

func TestThresholdNodeGlobalCutoff(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{50, 50, 50, 255})
	img.Set(1, 0, color.RGBA{200, 200, 200, 255})

	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &ThresholdNode{ID: "th1", Method: "global", Value: 128}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("ThresholdNode failed: %v", err)
	}

	r1, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	r2, _, _, _ := ctx.Images[0].At(1, 0).RGBA()
	if r1>>8 != 0 {
		t.Errorf("expected pixel below cutoff to become black, got %d", r1>>8)
	}
	if r2>>8 != 255 {
		t.Errorf("expected pixel above cutoff to become white, got %d", r2>>8)
	}
}

func TestThresholdNodeOtsuSeparatesBimodalHistogram(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 1))
	img.Set(0, 0, color.RGBA{20, 20, 20, 255})
	img.Set(1, 0, color.RGBA{30, 30, 30, 255})
	img.Set(2, 0, color.RGBA{220, 220, 220, 255})
	img.Set(3, 0, color.RGBA{230, 230, 230, 255})

	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &ThresholdNode{ID: "th2", Method: "otsu"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("ThresholdNode failed: %v", err)
	}

	darkVal, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	brightVal, _, _, _ := ctx.Images[0].At(3, 0).RGBA()
	if darkVal>>8 != 0 {
		t.Errorf("expected dark cluster to threshold to black, got %d", darkVal>>8)
	}
	if brightVal>>8 != 255 {
		t.Errorf("expected bright cluster to threshold to white, got %d", brightVal>>8)
	}
}

func TestThresholdNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &ThresholdNode{ID: "th3"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
