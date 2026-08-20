package nodes

import (
	"image"
	"testing"
)

func TestBandCompositeNodeMergesThreeBands(t *testing.T) {
	bandR := solidColorImage(2, 2, 200, 200, 200) // Bright band -> should map to R
	bandG := solidColorImage(2, 2, 50, 50, 50)    // Dim band -> should map to G
	bandB := solidColorImage(2, 2, 120, 120, 120) // Mid band -> should map to B

	ctx := &ProcessContext{Images: []image.Image{bandR, bandG, bandB}}
	node := &BandCompositeNode{ID: "composite1"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("BandCompositeNode failed: %v", err)
	}

	r, g, b, _ := ctx.Images[0].At(0, 0).RGBA()
	if r>>8 != 200 || g>>8 != 50 || b>>8 != 120 {
		t.Errorf("expected composite (200,50,120), got (%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func TestBandCompositeNodeRequiresThreeBands(t *testing.T) {
	ctx := &ProcessContext{Images: []image.Image{solidColorImage(2, 2, 1, 1, 1)}}
	node := &BandCompositeNode{ID: "composite2"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when fewer than three bands are provided")
	}
}
