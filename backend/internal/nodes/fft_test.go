package nodes

import (
	"image"
	"math"
	"math/cmplx"
	"testing"
)

func TestFFT1DConstantSignalHasOnlyDCEnergy(t *testing.T) {
	input := make([]complex128, 8)
	for i := range input {
		input[i] = complex(5, 0)
	}

	out := fft1D(input, false)

	if cmplx.Abs(out[0]-complex(40, 0)) > 1e-6 {
		t.Errorf("expected DC bin = 40, got %v", out[0])
	}
	for i := 1; i < len(out); i++ {
		if cmplx.Abs(out[i]) > 1e-6 {
			t.Errorf("expected bin %d to be ~0 for a constant signal, got %v", i, out[i])
		}
	}
}

func TestFFT1DRoundTripRecoversOriginalSignal(t *testing.T) {
	input := []complex128{1, 2, 3, 4, 5, 6, 7, 8}
	spectrum := fft1D(input, false)
	recovered := fft1D(spectrum, true)

	for i := range input {
		if cmplx.Abs(recovered[i]-input[i]) > 1e-6 {
			t.Errorf("expected round-trip FFT+IFFT to recover input at %d: got %v, want %v", i, recovered[i], input[i])
		}
	}
}

func TestNextPowerOfTwo(t *testing.T) {
	cases := map[int]int{1: 2, 2: 2, 3: 4, 5: 8, 8: 8, 9: 16}
	for in, want := range cases {
		if got := nextPowerOfTwo(in); got != want {
			t.Errorf("nextPowerOfTwo(%d) = %d, want %d", in, got, want)
		}
	}
}

func TestFFTNodeSpectrumModeOutputsPowerOfTwoSquare(t *testing.T) {
	img := solidColorImage(5, 3, 100, 100, 100)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &FFTNode{ID: "fft1", Mode: "spectrum"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("FFTNode failed: %v", err)
	}

	b := ctx.Images[0].Bounds()
	if b.Dx() != b.Dy() {
		t.Errorf("expected a square spectrum image, got %dx%d", b.Dx(), b.Dy())
	}
	size := b.Dx()
	if size&(size-1) != 0 {
		t.Errorf("expected spectrum size to be a power of two, got %d", size)
	}
}

func TestFFTNodeFilterModePreservesInputResolution(t *testing.T) {
	img := solidColorImage(6, 4, 128, 128, 128)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &FFTNode{ID: "fft2", Mode: "filter", Filter: "lowpass", CutoffRatio: 0.5}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("FFTNode failed: %v", err)
	}

	b := ctx.Images[0].Bounds()
	if b.Dx() != 6 || b.Dy() != 4 {
		t.Errorf("expected filtered output to match input resolution 6x4, got %dx%d", b.Dx(), b.Dy())
	}
}

func TestFFTNodeLowpassOnUniformImageStaysUniform(t *testing.T) {
	img := solidColorImage(8, 8, 100, 100, 100)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &FFTNode{ID: "fft3", Mode: "filter", Filter: "lowpass", CutoffRatio: 0.9}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("FFTNode failed: %v", err)
	}

	r, _, _, _ := ctx.Images[0].At(4, 4).RGBA()
	got := r >> 8
	if math.Abs(float64(got)-100) > 5 {
		t.Errorf("expected a wide lowpass filter to preserve a near-uniform image (~100), got %d", got)
	}
}

func TestFFTNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &FFTNode{ID: "fft4"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
