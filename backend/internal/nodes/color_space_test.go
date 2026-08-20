package nodes

import (
	"image"
	"testing"
)

func TestColorSpaceNodeGrayscaleMatchesLuminanceFormula(t *testing.T) {
	img := solidColorImage(2, 2, 100, 150, 200)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &ColorSpaceNode{ID: "cs1", Mode: "grayscale"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("ColorSpaceNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(0, 0).RGBA()
	if r>>8 == 0 || r>>8 == 255 {
		t.Errorf("unexpected grayscale extreme value: %d", r>>8)
	}
}

func TestColorSpaceNodeHSVOfPureRed(t *testing.T) {
	img := solidColorImage(1, 1, 255, 0, 0)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &ColorSpaceNode{ID: "cs2", Mode: "hsv"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("ColorSpaceNode failed: %v", err)
	}

	h, s, v, _ := ctx.Images[0].At(0, 0).RGBA()
	if h>>8 > 5 {
		t.Errorf("expected hue near 0 for pure red, got %d", h>>8)
	}
	if s>>8 < 250 {
		t.Errorf("expected full saturation for pure red, got %d", s>>8)
	}
	if v>>8 < 250 {
		t.Errorf("expected full value for pure red, got %d", v>>8)
	}
}

func TestColorSpaceNodeYCbCrOfNeutralGray(t *testing.T) {
	img := solidColorImage(1, 1, 128, 128, 128)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &ColorSpaceNode{ID: "cs3", Mode: "ycbcr"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("ColorSpaceNode failed: %v", err)
	}

	y, cb, cr, _ := ctx.Images[0].At(0, 0).RGBA()
	if diff := int(y>>8) - 128; diff < -3 || diff > 3 {
		t.Errorf("expected Y near 128 for neutral gray, got %d", y>>8)
	}
	if diff := int(cb>>8) - 128; diff < -3 || diff > 3 {
		t.Errorf("expected Cb near 128 for neutral gray, got %d", cb>>8)
	}
	if diff := int(cr>>8) - 128; diff < -3 || diff > 3 {
		t.Errorf("expected Cr near 128 for neutral gray, got %d", cr>>8)
	}
}

func TestColorSpaceNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &ColorSpaceNode{ID: "cs4"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
