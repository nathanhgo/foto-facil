package nodes

import (
	"image"
	"testing"
)

func TestHistogramNodeCountsChannels(t *testing.T) {
	img := solidColorImage(3, 2, 10, 20, 30) // 6 pixels, all identical
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &HistogramNode{ID: "hist1"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("HistogramNode failed: %v", err)
	}

	if node.Counts[0][10] != 6 {
		t.Errorf("expected 6 pixels at R=10, got %d", node.Counts[0][10])
	}
	if node.Counts[1][20] != 6 {
		t.Errorf("expected 6 pixels at G=20, got %d", node.Counts[1][20])
	}
	if node.Counts[2][30] != 6 {
		t.Errorf("expected 6 pixels at B=30, got %d", node.Counts[2][30])
	}

	if len(ctx.Images) != 1 {
		t.Fatalf("expected the node to output a single chart image, got %d", len(ctx.Images))
	}
}

func TestHistogramNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &HistogramNode{ID: "hist2"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
