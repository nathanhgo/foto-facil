package nodes

import (
	"image"
	"image/color"
	"testing"
)

func TestCropResizeNode_Resize(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	ctx := &ProcessContext{Images: []image.Image{img}}

	node := &CropResizeNode{ID: "crop_1", Width: 100, Height: 50}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("CropResizeNode failed: %v", err)
	}

	bounds := ctx.Images[0].Bounds()
	if bounds.Max.X != 100 || bounds.Max.Y != 50 {
		t.Errorf("Expected 100x50, got %dx%d", bounds.Max.X, bounds.Max.Y)
	}
}

func TestCropResizeNode_Crop(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	img.Set(50, 50, color.RGBA{255, 0, 0, 255})
	ctx := &ProcessContext{Images: []image.Image{img}}

	node := &CropResizeNode{ID: "crop_2", CropX: 40, CropY: 40, CropW: 30, CropH: 30}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("CropResizeNode failed: %v", err)
	}

	bounds := ctx.Images[0].Bounds()
	if bounds.Max.X != 30 || bounds.Max.Y != 30 {
		t.Errorf("Expected 30x30 crop, got %dx%d", bounds.Max.X, bounds.Max.Y)
	}
}

func TestCropResizeNode_NoInput(t *testing.T) {
	node := &CropResizeNode{ID: "crop_3"}
	if err := node.Process(&ProcessContext{}); err == nil {
		t.Fatal("Expected error for empty context")
	}
}

func TestRotateFlipNode_Rotate90(t *testing.T) {
	// 4x2 image
	img := image.NewRGBA(image.Rect(0, 0, 4, 2))
	ctx := &ProcessContext{Images: []image.Image{img}}

	node := &RotateFlipNode{ID: "rot_1", Angle: Rotate90, DoRotate: true}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("RotateFlipNode failed: %v", err)
	}

	b := ctx.Images[0].Bounds()
	// After 90° rotation of 4x2 → should be 2x4
	if b.Max.X != 2 || b.Max.Y != 4 {
		t.Errorf("Expected 2x4 after 90° rotation, got %dx%d", b.Max.X, b.Max.Y)
	}
}

func TestRotateFlipNode_FlipHorizontal(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255}) // Red at left
	ctx := &ProcessContext{Images: []image.Image{img}}

	node := &RotateFlipNode{ID: "flip_1", Flip: FlipHorizontal, DoFlip: true}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("RotateFlipNode failed: %v", err)
	}

	// Red pixel should now be at the right (x=3)
	r, _, _, _ := ctx.Images[0].At(3, 0).RGBA()
	if uint8(r>>8) != 255 {
		t.Errorf("Expected red at x=3 after horizontal flip, got %d", uint8(r>>8))
	}
}

func TestRotateFlipNode_NoInput(t *testing.T) {
	node := &RotateFlipNode{ID: "rot_2"}
	if err := node.Process(&ProcessContext{}); err == nil {
		t.Fatal("Expected error for empty context")
	}
}
