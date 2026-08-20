package nodes

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestStatisticsNodeComputesMeanMedianStdDev(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.Gray{Y: 10})
	img.Set(1, 0, color.Gray{Y: 20})
	img.Set(0, 1, color.Gray{Y: 30})
	img.Set(1, 1, color.Gray{Y: 40})

	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &StatisticsNode{ID: "stats1"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("StatisticsNode failed: %v", err)
	}

	if node.Mean != 25 {
		t.Errorf("expected mean 25, got %f", node.Mean)
	}
	if node.Median != 25 {
		t.Errorf("expected median 25, got %f", node.Median)
	}
	if node.StdDev <= 0 {
		t.Errorf("expected a positive standard deviation, got %f", node.StdDev)
	}
}

func TestStatisticsNodeComputesPSNRAgainstOriginal(t *testing.T) {
	img := solidColorImage(2, 2, 100, 100, 100)
	ctx := &ProcessContext{
		Images:         []image.Image{img},
		OriginalImages: []image.Image{img},
	}
	node := &StatisticsNode{ID: "stats2"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("StatisticsNode failed: %v", err)
	}

	if node.MSE != 0 {
		t.Errorf("expected MSE 0 for identical images, got %f", node.MSE)
	}
	if !math.IsInf(node.PSNR, 1) {
		t.Errorf("expected PSNR to be +Inf for identical images, got %f", node.PSNR)
	}
}

func TestStatisticsNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &StatisticsNode{ID: "stats3"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
