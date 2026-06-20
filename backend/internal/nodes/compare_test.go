package nodes

import (
	"image"
	"image/color"
	"testing"
)

func TestCompareNode(t *testing.T) {
	ctx := &ProcessContext{}

	// Two 10x10 images: red and blue
	imgA := image.NewRGBA(image.Rect(0, 0, 10, 10))
	imgB := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			imgA.Set(x, y, color.RGBA{255, 0, 0, 255})
			imgB.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}
	ctx.OriginalImages = []image.Image{imgA}
	ctx.Images = []image.Image{imgB}

	node := &CompareNode{ID: "compare_1"}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("CompareNode failed: %v", err)
	}

	// Should have replaced with the combined image
	if len(ctx.Images) != 1 {
		t.Fatalf("Expected exactly 1 combined image, got %d images", len(ctx.Images))
	}

	combined := ctx.Images[0]
	bounds := combined.Bounds()
	// Width should be 10 + 2 (divider) + 10 = 22
	if bounds.Max.X != 22 {
		t.Errorf("Expected combined width of 22, got %d", bounds.Max.X)
	}
}

func TestCompareNode_InsufficientInput(t *testing.T) {
	ctx := &ProcessContext{Images: []image.Image{image.NewRGBA(image.Rect(0, 0, 5, 5))}} // no OriginalImages
	node := &CompareNode{ID: "compare_2"}
	if err := node.Process(ctx); err == nil {
		t.Fatal("Expected error with no OriginalImages")
	}
}
