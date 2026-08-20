package nodes

import (
	"errors"
	"image"
	"image/color"
	"math"
	"math/cmplx"
)

// FFTNode computes the 2D Fast Fourier Transform of an image's luminance
// channel. In "spectrum" mode (default) it renders the magnitude spectrum,
// a staple visualization in signal processing courses — as idea.md puts it,
// "o terror e o amor de todo estudante de PDI". In "filter" mode it applies
// an ideal low-pass or high-pass filter in the frequency domain and
// transforms back to the spatial domain.
type FFTNode struct {
	ID          string
	Mode        string  // "spectrum" (default) or "filter"
	Filter      string  // "lowpass" (default) or "highpass", used when Mode == "filter"
	CutoffRatio float64 // Fraction (0-1) of the spectrum radius kept/removed; default 0.2
}

func (n *FFTNode) GetID() string { return n.ID }

func (n *FFTNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("FFTNode requires at least one input image")
	}

	mode := n.Mode
	if mode == "" {
		mode = "spectrum"
	}
	cutoff := n.CutoffRatio
	if cutoff <= 0 {
		cutoff = 0.2
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()
		size := nextPowerOfTwo(imax(w, h))

		grid := make([][]complex128, size)
		for y := 0; y < size; y++ {
			grid[y] = make([]complex128, size)
			if y >= h {
				continue
			}
			for x := 0; x < w; x++ {
				r, g, bch, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
				grid[y][x] = complex(float64(luminance8(r, g, bch)), 0)
			}
		}

		spectrum := fft2D(grid, size, false)

		if mode == "filter" {
			applyFrequencyFilter(spectrum, size, n.Filter, cutoff)
			spatial := fft2D(spectrum, size, true)
			dst := image.NewGray(image.Rect(0, 0, w, h))
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					dst.Set(x, y, color.Gray{Y: uint8(clampFloat(real(spatial[y][x])))})
				}
			}
			processed = append(processed, dst)
			continue
		}

		processed = append(processed, renderSpectrum(spectrum, size))
	}

	ctx.Images = processed
	return nil
}

// renderSpectrum converts a raw 2D FFT result into a human-readable,
// log-scaled magnitude image with the DC (zero-frequency) term centered.
func renderSpectrum(spectrum [][]complex128, size int) image.Image {
	mags := make([][]float64, size)
	maxMag := 0.0
	for y := 0; y < size; y++ {
		mags[y] = make([]float64, size)
		for x := 0; x < size; x++ {
			mag := math.Log(1 + cmplx.Abs(spectrum[y][x]))
			mags[y][x] = mag
			if mag > maxMag {
				maxMag = mag
			}
		}
	}
	if maxMag == 0 {
		maxMag = 1
	}

	dst := image.NewGray(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sx := (x + size/2) % size
			sy := (y + size/2) % size
			v := mags[sy][sx] / maxMag * 255
			dst.Set(x, y, color.Gray{Y: uint8(clampFloat(v))})
		}
	}
	return dst
}

func nextPowerOfTwo(n int) int {
	p := 2
	for p < n {
		p *= 2
	}
	return p
}

// fft2D runs a separable 2D FFT (rows then columns). Set inverse to true to
// run the inverse transform (IFFT), which also normalizes the output.
func fft2D(grid [][]complex128, size int, inverse bool) [][]complex128 {
	result := make([][]complex128, size)
	for y := 0; y < size; y++ {
		result[y] = fft1D(grid[y], inverse)
	}
	for x := 0; x < size; x++ {
		col := make([]complex128, size)
		for y := 0; y < size; y++ {
			col[y] = result[y][x]
		}
		col = fft1D(col, inverse)
		for y := 0; y < size; y++ {
			result[y][x] = col[y]
		}
	}
	return result
}

// fft1D is an iterative radix-2 Cooley-Tukey FFT. len(input) must be a power
// of two.
func fft1D(input []complex128, inverse bool) []complex128 {
	n := len(input)
	a := make([]complex128, n)
	copy(a, input)

	// Bit-reversal permutation
	for i, j := 1, 0; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j ^= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}

	sign := -1.0
	if inverse {
		sign = 1.0
	}

	for length := 2; length <= n; length <<= 1 {
		angle := sign * 2 * math.Pi / float64(length)
		wLen := cmplx.Rect(1, angle)
		half := length / 2
		for start := 0; start < n; start += length {
			w := complex(1.0, 0.0)
			for k := 0; k < half; k++ {
				u := a[start+k]
				v := a[start+k+half] * w
				a[start+k] = u + v
				a[start+k+half] = u - v
				w *= wLen
			}
		}
	}

	if inverse {
		for i := range a {
			a[i] /= complex(float64(n), 0)
		}
	}

	return a
}

// applyFrequencyFilter zeroes out spectrum components based on their
// distance from the DC term, keeping only low frequencies ("lowpass") or
// only high frequencies ("highpass"). No fftshift is needed since distance
// is computed with wrap-around directly on the natural FFT layout.
func applyFrequencyFilter(spectrum [][]complex128, size int, filter string, cutoffRatio float64) {
	if filter == "" {
		filter = "lowpass"
	}
	cutoff := cutoffRatio * float64(size) / 2

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(imin(x, size-x))
			dy := float64(imin(y, size-y))
			dist := math.Sqrt(dx*dx + dy*dy)

			if filter == "highpass" {
				if dist < cutoff {
					spectrum[y][x] = 0
				}
			} else if dist > cutoff {
				spectrum[y][x] = 0
			}
		}
	}
}
