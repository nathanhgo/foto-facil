package nodes

import (
	"image"
	"image/color"
	"testing"
)

// markedImage builds a mostly-dark image with a single bright marker pixel,
// useful to verify that alignment recovers a known translation.
func markedImage(size, markerX, markerY int) *image.RGBA {
	img := solidColorImage(size, size, 10, 10, 10)
	img.Set(markerX, markerY, color.RGBA{250, 250, 250, 255})
	return img
}

func TestRegistrationNodeAlignsShiftedImage(t *testing.T) {
	reference := markedImage(20, 10, 10)
	moving := markedImage(20, 13, 8) // Marker shifted by (+3, -2) relative to reference

	ctx := &ProcessContext{Images: []image.Image{reference, moving}}
	node := &RegistrationNode{ID: "reg1", MaxShift: 6}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("RegistrationNode failed: %v", err)
	}
	if len(ctx.Images) != 2 {
		t.Fatalf("expected reference + aligned image, got %d images", len(ctx.Images))
	}

	aligned := ctx.Images[1]
	r, g, b, _ := aligned.At(10, 10).RGBA()
	if r>>8 < 200 || g>>8 < 200 || b>>8 < 200 {
		t.Errorf("expected the marker to be realigned to (10,10), got brightness (%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func TestRegistrationNodeRequiresTwoImages(t *testing.T) {
	ctx := &ProcessContext{Images: []image.Image{solidColorImage(4, 4, 1, 1, 1)}}
	node := &RegistrationNode{ID: "reg2"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when fewer than two images are provided")
	}
}
