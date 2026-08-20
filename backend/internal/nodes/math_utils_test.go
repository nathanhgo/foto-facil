package nodes

import (
	"image"
	"image/color"
	"testing"
)

// solidColorImage builds a w x h RGBA test image filled with a single color.
// Shared across the node test files in this package.
func solidColorImage(w, h int, r, g, b uint8) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

func TestImaxImin(t *testing.T) {
	if imax(3, 7) != 7 {
		t.Errorf("expected imax(3,7) = 7")
	}
	if imin(3, 7) != 3 {
		t.Errorf("expected imin(3,7) = 3")
	}
}

func TestClampInt(t *testing.T) {
	if clampInt(-5, 0, 255) != 0 {
		t.Errorf("expected clampInt(-5, 0, 255) = 0")
	}
	if clampInt(300, 0, 255) != 255 {
		t.Errorf("expected clampInt(300, 0, 255) = 255")
	}
	if clampInt(100, 0, 255) != 100 {
		t.Errorf("expected clampInt(100, 0, 255) = 100")
	}
}

func TestClampFloat(t *testing.T) {
	if clampFloat(-10) != 0 {
		t.Errorf("expected clampFloat(-10) = 0")
	}
	if clampFloat(500) != 255 {
		t.Errorf("expected clampFloat(500) = 255")
	}
	if clampFloat(128) != 128 {
		t.Errorf("expected clampFloat(128) = 128")
	}
}

func TestLuminance8(t *testing.T) {
	// Pure white at 16-bit full scale should luminance to 255
	if got := luminance8(0xFFFF, 0xFFFF, 0xFFFF); got != 255 {
		t.Errorf("expected luminance8(white) = 255, got %d", got)
	}
	// Pure black should luminance to 0
	if got := luminance8(0, 0, 0); got != 0 {
		t.Errorf("expected luminance8(black) = 0, got %d", got)
	}
}

func TestGrayBufferToGrayValuesAndAt(t *testing.T) {
	img := solidColorImage(2, 2, 100, 150, 200)
	g := toGrayValues(img)

	if g.w != 2 || g.h != 2 {
		t.Fatalf("expected 2x2 gray buffer, got %dx%d", g.w, g.h)
	}

	v, ok := g.at(0, 0)
	if !ok {
		t.Fatal("expected (0,0) to be in bounds")
	}
	r, gg, b, _ := color.RGBA{100, 150, 200, 255}.RGBA()
	expected := luminance8(r, gg, b)
	if v != expected {
		t.Errorf("expected gray value %d, got %d", expected, v)
	}

	if _, ok := g.at(-1, 0); ok {
		t.Error("expected out-of-bounds access to return ok=false")
	}
}
