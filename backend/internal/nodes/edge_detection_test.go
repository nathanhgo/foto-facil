package nodes

import (
	"image"
	"image/color"
	"testing"
)

// verticalEdgeImage builds an image whose left half is black and right half
// is white, producing a strong vertical edge in the middle column.
func verticalEdgeImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x < w/2 {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}
	return img
}

func TestEdgeDetectionNodeSobelHighlightsEdge(t *testing.T) {
	img := verticalEdgeImage(6, 4)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &EdgeDetectionNode{ID: "edge1", Method: "sobel"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("EdgeDetectionNode failed: %v", err)
	}

	edgeVal, _, _, _ := ctx.Images[0].At(3, 2).RGBA()
	flatVal, _, _, _ := ctx.Images[0].At(0, 2).RGBA()

	if edgeVal>>8 <= flatVal>>8 {
		t.Errorf("expected the edge column (x=3) to have a stronger response than the flat region (x=0): got %d vs %d", edgeVal>>8, flatVal>>8)
	}
}

func TestEdgeDetectionNodeLaplacianHighlightsEdge(t *testing.T) {
	img := verticalEdgeImage(6, 4)
	ctx := &ProcessContext{Images: []image.Image{img}}
	node := &EdgeDetectionNode{ID: "edge2", Method: "laplacian"}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("EdgeDetectionNode failed: %v", err)
	}

	edgeVal, _, _, _ := ctx.Images[0].At(3, 2).RGBA()
	flatVal, _, _, _ := ctx.Images[0].At(0, 2).RGBA()

	if edgeVal>>8 <= flatVal>>8 {
		t.Errorf("expected the edge column to have a stronger Laplacian response than the flat region: got %d vs %d", edgeVal>>8, flatVal>>8)
	}
}

func TestEdgeDetectionNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &EdgeDetectionNode{ID: "edge3"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
