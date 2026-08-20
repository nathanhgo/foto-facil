package nodes

import (
	"errors"
	"image"
	"math"
	"sort"
)

// StatisticsNode computes summary metrics for the input image — mean, median
// and standard deviation of pixel luminance — and, when an original
// reference is available in the pipeline (ctx.OriginalImages), similarity
// metrics against it (MSE and PSNR). The image itself passes through
// unchanged; surfacing these numbers in the properties panel UI is a
// follow-up (see mvp.md) — for now they are exposed on the struct fields.
type StatisticsNode struct {
	ID string

	Mean   float64
	Median float64
	StdDev float64
	MSE    float64
	PSNR   float64
}

func (n *StatisticsNode) GetID() string { return n.ID }

func (n *StatisticsNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("StatisticsNode requires at least one input image")
	}

	values := luminanceValues(ctx.Images[0])

	n.Mean = meanOf(values)
	n.Median = medianOf(values)
	n.StdDev = stdDevOf(values, n.Mean)

	if len(ctx.OriginalImages) > 0 {
		refValues := luminanceValues(ctx.OriginalImages[0])
		n.MSE = meanSquaredError(values, refValues)
		n.PSNR = psnrOf(n.MSE)
	}

	return nil
}

func luminanceValues(img image.Image) []float64 {
	b := img.Bounds()
	values := make([]float64, 0, b.Dx()*b.Dy())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bch, _ := img.At(x, y).RGBA()
			values = append(values, float64(luminance8(r, g, bch)))
		}
	}
	return values
}

func meanOf(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func medianOf(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

func stdDevOf(values []float64, m float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sumSq float64
	for _, v := range values {
		sumSq += (v - m) * (v - m)
	}
	return math.Sqrt(sumSq / float64(len(values)))
}

func meanSquaredError(a, b []float64) float64 {
	n := imin(len(a), len(b))
	if n == 0 {
		return 0
	}
	var sum float64
	for i := 0; i < n; i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum / float64(n)
}

func psnrOf(mse float64) float64 {
	if mse == 0 {
		return math.Inf(1)
	}
	return 20*math.Log10(255) - 10*math.Log10(mse)
}
