package nodes

import (
	"image"
	"testing"
)

func TestPixelMatrixNodeRendersRegionGrid(t *testing.T) {
	img := solidColorImage(10, 10, 128, 128, 128)
	ctx := &ProcessContext{Images: []image.Image{img}}

	node := &PixelMatrixNode{ID: "matrix1", RegionSize: 4}
	if err := node.Process(ctx); err != nil {
		t.Fatalf("PixelMatrixNode failed: %v", err)
	}

	out := ctx.Images[0]
	b := out.Bounds()
	expectedSize := 4 * pixelMatrixCellSize
	if b.Dx() != expectedSize || b.Dy() != expectedSize {
		t.Errorf("expected canvas %dx%d, got %dx%d", expectedSize, expectedSize, b.Dx(), b.Dy())
	}

	// Sample a pixel deep inside the first cell (away from the text) to
	// confirm the cell background reflects the source luminance (~128).
	r, _, _, _ := out.At(pixelMatrixCellSize-2, pixelMatrixCellSize-2).RGBA()
	if r>>8 < 120 || r>>8 > 136 {
		t.Errorf("expected cell background shaded near luminance 128, got %d", r>>8)
	}
}

func TestPixelMatrixNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &PixelMatrixNode{ID: "matrix2"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
