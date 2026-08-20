package nodes

import (
	"image"
	"testing"
)

func TestAugmentationNodeProducesOneOutputPerInput(t *testing.T) {
	img1 := solidColorImage(6, 6, 100, 120, 140)
	img2 := solidColorImage(6, 6, 100, 120, 140)

	ctx := &ProcessContext{Images: []image.Image{img1, img2}}
	node := &AugmentationNode{ID: "aug1", Seed: 10}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("AugmentationNode failed: %v", err)
	}
	if len(ctx.Images) != 2 {
		t.Fatalf("expected 2 augmented outputs, got %d", len(ctx.Images))
	}
}

func TestAugmentationNodeIsDeterministicGivenSameSeed(t *testing.T) {
	img1 := solidColorImage(6, 6, 100, 120, 140)
	img2 := solidColorImage(6, 6, 100, 120, 140)

	ctx1 := &ProcessContext{Images: []image.Image{img1}}
	ctx2 := &ProcessContext{Images: []image.Image{img2}}

	node1 := &AugmentationNode{ID: "aug2", Seed: 99}
	node2 := &AugmentationNode{ID: "aug3", Seed: 99}

	if err := node1.Process(ctx1); err != nil {
		t.Fatalf("AugmentationNode failed: %v", err)
	}
	if err := node2.Process(ctx2); err != nil {
		t.Fatalf("AugmentationNode failed: %v", err)
	}

	b1 := ctx1.Images[0].Bounds()
	b2 := ctx2.Images[0].Bounds()
	if b1 != b2 {
		t.Errorf("expected the same seed to produce the same transformation (bounds), got %v vs %v", b1, b2)
	}

	r1, _, _, _ := ctx1.Images[0].At(0, 0).RGBA()
	r2, _, _, _ := ctx2.Images[0].At(0, 0).RGBA()
	if r1 != r2 {
		t.Error("expected the same seed to produce deterministic, reproducible augmentation")
	}
}

func TestAugmentationNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &AugmentationNode{ID: "aug4"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
