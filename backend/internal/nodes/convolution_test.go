package nodes

import (
	"image"
	"testing"
)

func TestConvolutionNodeIdentityKernelPreservesImage(t *testing.T) {
	img := solidColorImage(3, 3, 77, 88, 99)
	ctx := &ProcessContext{Images: []image.Image{img}}

	node := &ConvolutionNode{
		ID:         "conv1",
		KernelSize: 3,
		Kernel:     []float64{0, 0, 0, 0, 1, 0, 0, 0, 0},
	}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("ConvolutionNode failed: %v", err)
	}

	r, g, b, _ := ctx.Images[0].At(1, 1).RGBA()
	if r>>8 != 77 || g>>8 != 88 || b>>8 != 99 {
		t.Errorf("expected identity kernel to preserve pixel (77,88,99), got (%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func TestConvolutionNodeBoxBlurOnUniformImageStaysUniform(t *testing.T) {
	img := solidColorImage(4, 4, 150, 150, 150)
	ctx := &ProcessContext{Images: []image.Image{img}}

	node := &ConvolutionNode{
		ID:         "conv2",
		KernelSize: 3,
		Kernel:     []float64{1, 1, 1, 1, 1, 1, 1, 1, 1}, // Averaging kernel (normalized internally)
	}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("ConvolutionNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(2, 2).RGBA()
	if r>>8 != 150 {
		t.Errorf("expected uniform image to remain uniform after box blur, got %d", r>>8)
	}
}

func TestConvolutionNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &ConvolutionNode{ID: "conv3"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
