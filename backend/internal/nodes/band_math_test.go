package nodes

import (
	"image"
	"testing"
)

func TestBandMathNodeNormalizedDifference(t *testing.T) {
	// A = 200 (bright, e.g. simulated near-infrared reflectance)
	// B = 50 (dim, e.g. red band)
	// NDVI-like index = (200-50)/(200+50) = 0.6 -> scaled to (0.6+1)/2*255 ≈ 204
	bandA := solidColorImage(2, 2, 200, 200, 200)
	bandB := solidColorImage(2, 2, 50, 50, 50)

	ctx := &ProcessContext{Images: []image.Image{bandA, bandB}}
	node := &BandMathNode{ID: "math1"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("BandMathNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	got := int(r >> 8)
	if got < 200 || got > 210 {
		t.Errorf("expected normalized difference around 204, got %d", got)
	}
}

func TestBandMathNodeEqualBandsYieldsMidpoint(t *testing.T) {
	bandA := solidColorImage(1, 1, 100, 100, 100)
	bandB := solidColorImage(1, 1, 100, 100, 100)

	ctx := &ProcessContext{Images: []image.Image{bandA, bandB}}
	node := &BandMathNode{ID: "math2"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("BandMathNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	got := int(r >> 8)
	if got < 125 || got > 130 {
		t.Errorf("expected equal bands to yield a mid-gray index (~127), got %d", got)
	}
}

func TestBandMathNodeRequiresTwoBands(t *testing.T) {
	ctx := &ProcessContext{Images: []image.Image{solidColorImage(2, 2, 1, 1, 1)}}
	node := &BandMathNode{ID: "math3"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when fewer than two bands are provided")
	}
}
