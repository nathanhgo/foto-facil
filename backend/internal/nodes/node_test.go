package nodes

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func TestInputNode(t *testing.T) {
	// Create a real temporary JPEG image to test actual file reading
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "test.jpg")

	srcImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			srcImg.Set(x, y, color.RGBA{100, 150, 200, 255})
		}
	}
	f, _ := os.Create(imgPath)
	jpeg.Encode(f, srcImg, nil)
	f.Close()

	node := &InputNode{ID: "node_input_1", FilePaths: []string{imgPath}}
	ctx := &ProcessContext{}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("InputNode failed: %v", err)
	}
	if len(ctx.Images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(ctx.Images))
	}
}

func TestBrightnessNode(t *testing.T) {
	ctx := &ProcessContext{}
	
	// Create a dummy image (1x1) with color RGB(100, 100, 100)
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{100, 100, 100, 255})
	ctx.Images = append(ctx.Images, img)

	// Apply +50 brightness
	node := &BrightnessNode{
		ID:         "node_bright_1",
		Brightness: 50,
	}

	err := node.Process(ctx)
	if err != nil {
		t.Fatalf("BrightnessNode failed: %v", err)
	}

	processedImg := ctx.Images[0]
	r, g, b, _ := processedImg.At(0, 0).RGBA()
	
	// Convert from 16-bit to 8-bit for assertion
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

	if r8 != 150 || g8 != 150 || b8 != 150 {
		t.Errorf("Expected pixel color to be (150, 150, 150), got (%d, %d, %d)", r8, g8, b8)
	}

	// Test Clamp (ex: +200 should clamp at 255)
	node.Brightness = 200
	node.Process(ctx) // Process the already 150-brightness image
	
	r, _, _, _ = ctx.Images[0].At(0, 0).RGBA()
	r8 = uint8(r >> 8)
	if r8 != 255 {
		t.Errorf("Expected pixel to be clamped at 255, got %d", r8)
	}
}
