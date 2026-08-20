package nodes

import (
	"image"
	"testing"
)

func TestNoiseNodeGaussianAltersPixelsButStaysInRange(t *testing.T) {
	img := solidColorImage(6, 6, 128, 128, 128)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &NoiseNode{ID: "noise1", Type: "gaussian", Amount: 30, Seed: 7}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("NoiseNode failed: %v", err)
	}

	out := ctx.Images[0]
	changed := false
	b := out.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, _, _, _ := out.At(x, y).RGBA()
			v := r >> 8
			if v > 255 {
				t.Fatalf("pixel value out of range: %d", v)
			}
			if v != 128 {
				changed = true
			}
		}
	}
	if !changed {
		t.Error("expected gaussian noise to alter at least some pixels")
	}
}

func TestNoiseNodeSaltPepperProducesExtremes(t *testing.T) {
	img := solidColorImage(20, 20, 128, 128, 128)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &NoiseNode{ID: "noise2", Type: "salt_pepper", Amount: 0.5, Seed: 3}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("NoiseNode failed: %v", err)
	}

	out := ctx.Images[0]
	foundExtreme := false
	b := out.Bounds()
	for y := b.Min.Y; y < b.Max.Y && !foundExtreme; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, _, _, _ := out.At(x, y).RGBA()
			v := r >> 8
			if v == 0 || v == 255 {
				foundExtreme = true
				break
			}
		}
	}
	if !foundExtreme {
		t.Error("expected salt & pepper noise to produce at least one pure black or white pixel")
	}
}

func TestNoiseNodeIsDeterministicGivenSameSeed(t *testing.T) {
	img1 := solidColorImage(4, 4, 100, 100, 100)
	img2 := solidColorImage(4, 4, 100, 100, 100)

	ctx1 := &ProcessContext{Images: []image.Image{img1}}
	ctx2 := &ProcessContext{Images: []image.Image{img2}}

	node1 := &NoiseNode{ID: "noise3", Type: "gaussian", Amount: 20, Seed: 42}
	node2 := &NoiseNode{ID: "noise4", Type: "gaussian", Amount: 20, Seed: 42}

	if err := node1.Process(ctx1); err != nil {
		t.Fatalf("NoiseNode failed: %v", err)
	}
	if err := node2.Process(ctx2); err != nil {
		t.Fatalf("NoiseNode failed: %v", err)
	}

	r1, _, _, _ := ctx1.Images[0].At(1, 1).RGBA()
	r2, _, _, _ := ctx2.Images[0].At(1, 1).RGBA()
	if r1 != r2 {
		t.Error("expected the same seed to produce deterministic, reproducible noise")
	}
}

func TestNoiseNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &NoiseNode{ID: "noise5"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
