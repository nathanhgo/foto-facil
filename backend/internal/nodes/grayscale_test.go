package nodes

import (
	"image"
	"image/color"
	"testing"
)

func TestGrayscaleNode(t *testing.T) {
	ctx := &ProcessContext{}

	// Pixel vermelho puro
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	ctx.Images = append(ctx.Images, img)

	node := &GrayscaleNode{ID: "node_gray_1"}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("GrayscaleNode failed: %v", err)
	}

	// Pixel vermelho (255,0,0) via luminosidade → ~76
	r, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	lum := uint8(r >> 8)
	if lum < 70 || lum > 82 {
		t.Errorf("Expected grayscale luminosity ~76 for red pixel, got %d", lum)
	}
}

func TestGrayscaleNodeNoInput(t *testing.T) {
	node := &GrayscaleNode{ID: "node_gray_2"}
	err := node.Process(&ProcessContext{})
	if err == nil {
		t.Fatal("Expected error when no input images provided")
	}
}
