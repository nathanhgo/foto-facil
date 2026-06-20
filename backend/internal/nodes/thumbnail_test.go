package nodes

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestThumbnailNode(t *testing.T) {
	ctx := &ProcessContext{}

	// Imagem 400x400 para testar redução
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	for y := 0; y < 400; y++ {
		for x := 0; x < 400; x++ {
			img.Set(x, y, color.RGBA{100, 150, 200, 255})
		}
	}
	ctx.Images = append(ctx.Images, img)

	node := &ThumbnailNode{ID: "thumb_1", MaxWidth: 100, MaxHeight: 100}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("ThumbnailNode failed: %v", err)
	}

	if len(node.Thumbnails) != 1 {
		t.Fatalf("Expected 1 thumbnail, got %d", len(node.Thumbnails))
	}

	thumb := node.Thumbnails[0]
	if !strings.HasPrefix(thumb, "data:image/png;base64,") {
		t.Errorf("Expected base64 PNG data URL, got: %s...", thumb[:30])
	}
}

func TestThumbnailNodeNoInput(t *testing.T) {
	node := &ThumbnailNode{ID: "thumb_2"}
	err := node.Process(&ProcessContext{})
	if err == nil {
		t.Fatal("Expected error when no input images provided")
	}
}
