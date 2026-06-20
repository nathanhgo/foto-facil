package nodes

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestOutputNode(t *testing.T) {
	ctx := &ProcessContext{}
	
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	ctx.Images = append(ctx.Images, img)

	tempDir := t.TempDir()

	node := &OutputNode{
		ID:        "node_out_1",
		OutputDir: tempDir,
	}

	err := node.Process(ctx)
	if err != nil {
		t.Fatalf("OutputNode failed: %v", err)
	}

	expectedFile := filepath.Join(tempDir, "output_node_out_1_0.png")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Expected output file was not created: %s", expectedFile)
	}
}
