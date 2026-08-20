package nodes

import (
	"image"
	"image/color"
	"testing"
)

func TestBlurNodeMeanOnUniformImageStaysUniform(t *testing.T) {
	img := solidColorImage(5, 5, 100, 100, 100)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &BlurNode{ID: "blur1", Method: "mean", KernelSize: 3}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("BlurNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(2, 2).RGBA()
	if r>>8 != 100 {
		t.Errorf("expected uniform image to remain uniform after mean blur, got %d", r>>8)
	}
}

func TestBlurNodeMedianRemovesSaltNoise(t *testing.T) {
	img := solidColorImage(5, 5, 50, 50, 50)
	// Inject a single "salt" outlier pixel in the middle.
	img.Set(2, 2, color.RGBA{255, 255, 255, 255})

	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &BlurNode{ID: "blur2", Method: "median", KernelSize: 3}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("BlurNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(2, 2).RGBA()
	if r>>8 != 50 {
		t.Errorf("expected median filter to remove the salt-noise outlier, got %d", r>>8)
	}
}

func TestBlurNodeGaussianSmoothsSharpEdge(t *testing.T) {
	img := verticalEdgeImage(10, 4)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &BlurNode{ID: "blur3", Method: "gaussian", KernelSize: 5}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("BlurNode failed: %v", err)
	}

	// The pixel right at the edge column should no longer be pure black or
	// pure white after Gaussian smoothing.
	r, _, _, _ := ctx.Images[0].At(5, 2).RGBA()
	if r>>8 == 0 || r>>8 == 255 {
		t.Errorf("expected the edge to be smoothed to an intermediate value, got %d", r>>8)
	}
}

func TestBlurNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &BlurNode{ID: "blur4"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
