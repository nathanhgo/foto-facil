package nodes

import (
	"image"
	"image/color"
	"testing"
)

func TestAINode_BackgroundRemoval(t *testing.T) {
	ctx := &ProcessContext{}
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	
	// Fill background with green screen, middle pixel with red
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{0, 255, 0, 255})
		}
	}
	img.Set(5, 5, color.RGBA{255, 0, 0, 255})
	ctx.Images = append(ctx.Images, img)

	node := &AINode{ID: "ai_1", Tool: "background_removal"}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("AINode failed: %v", err)
	}

	// Top-left pixel should now be transparent
	_, _, _, a := ctx.Images[0].At(0, 0).RGBA()
	if uint8(a>>8) != 0 {
		t.Errorf("Expected transparent background, got alpha %d", uint8(a>>8))
	}

	// Middle pixel should remain opaque
	_, _, _, a2 := ctx.Images[0].At(5, 5).RGBA()
	if uint8(a2>>8) != 255 {
		t.Errorf("Expected foreground pixel to remain opaque, got alpha %d", uint8(a2>>8))
	}
}

func TestAINode_ContrastBoost(t *testing.T) {
	ctx := &ProcessContext{}
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{128, 128, 128, 255}) // mid tone
	ctx.Images = append(ctx.Images, img)

	node := &AINode{ID: "ai_2", Tool: "contrast_boost"}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("AINode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	r8 := uint8(r >> 8)
	// 128 (0.5) under sigmoid 1 / (1 + exp(-10*(x-0.5))) -> 0.5 -> 127
	if r8 < 120 || r8 > 135 {
		t.Errorf("Expected boosted contrast to keep middle close, got %d", r8)
	}
}

func TestAINode_NoInput(t *testing.T) {
	node := &AINode{ID: "ai_3"}
	if err := node.Process(&ProcessContext{}); err == nil {
		t.Fatal("Expected error with empty context")
	}
}
