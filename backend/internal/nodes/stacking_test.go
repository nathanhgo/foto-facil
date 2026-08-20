package nodes

import (
	"image"
	"testing"
)

func TestStackingNodeAveragesFrames(t *testing.T) {
	frame1 := solidColorImage(2, 2, 100, 100, 100)
	frame2 := solidColorImage(2, 2, 200, 200, 200)

	ctx := &ProcessContext{Images: []image.Image{frame1, frame2}}
	node := &StackingNode{ID: "stack1"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("StackingNode failed: %v", err)
	}
	if len(ctx.Images) != 1 {
		t.Fatalf("expected a single averaged output image, got %d", len(ctx.Images))
	}

	r, g, b, _ := ctx.Images[0].At(0, 0).RGBA()
	r8, g8, b8 := r>>8, g>>8, b>>8
	if r8 != 150 || g8 != 150 || b8 != 150 {
		t.Errorf("expected averaged pixel (150,150,150), got (%d,%d,%d)", r8, g8, b8)
	}
}

func TestStackingNodeSkipsMismatchedFrames(t *testing.T) {
	frame1 := solidColorImage(2, 2, 100, 100, 100)
	frame2 := solidColorImage(2, 2, 200, 200, 200)
	mismatched := solidColorImage(5, 5, 0, 0, 0)

	ctx := &ProcessContext{Images: []image.Image{frame1, frame2, mismatched}}
	node := &StackingNode{ID: "stack2"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("StackingNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	if r>>8 != 150 {
		t.Errorf("expected mismatched-size frame to be skipped, still averaging to 150, got %d", r>>8)
	}
}

func TestStackingNodeRequiresAtLeastTwoFrames(t *testing.T) {
	ctx := &ProcessContext{Images: []image.Image{solidColorImage(2, 2, 1, 1, 1)}}
	node := &StackingNode{ID: "stack3"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when fewer than two frames are provided")
	}
}
