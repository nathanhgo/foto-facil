package nodes

import (
	"errors"
	"image"
	"image/color"
)

const histogramBins = 256

// HistogramNode computes the per-channel intensity distribution of an image
// and renders it as a bar chart — the most fundamental visualization tool in
// digital image processing.
type HistogramNode struct {
	ID string
	// Counts exposes the last computed per-channel histogram ([0]=R, [1]=G,
	// [2]=B) so other code (tests, future nodes) can inspect the raw numbers
	// instead of parsing pixels back out of the rendered chart image.
	Counts [3][histogramBins]int
}

func (n *HistogramNode) GetID() string { return n.ID }

func (n *HistogramNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("HistogramNode requires at least one input image")
	}

	img := ctx.Images[0]
	b := img.Bounds()

	var counts [3][histogramBins]int
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bch, _ := img.At(x, y).RGBA()
			counts[0][r>>8]++
			counts[1][g>>8]++
			counts[2][bch>>8]++
		}
	}
	n.Counts = counts

	const chartW, chartH = histogramBins, 120
	chart := image.NewRGBA(image.Rect(0, 0, chartW, chartH))
	for y := 0; y < chartH; y++ {
		for x := 0; x < chartW; x++ {
			chart.Set(x, y, color.RGBA{20, 20, 20, 255}) // Near-black background
		}
	}

	maxCount := 1
	for c := 0; c < 3; c++ {
		for _, v := range counts[c] {
			if v > maxCount {
				maxCount = v
			}
		}
	}

	channelColors := [3]color.RGBA{
		{255, 80, 80, 255},
		{80, 255, 80, 255},
		{80, 80, 255, 255},
	}

	for c := 0; c < 3; c++ {
		for x := 0; x < chartW; x++ {
			barHeight := counts[c][x] * (chartH - 1) / maxCount
			for y := 0; y < barHeight; y++ {
				chart.Set(x, chartH-1-y, channelColors[c])
			}
		}
	}

	ctx.Images = []image.Image{chart}
	return nil
}
