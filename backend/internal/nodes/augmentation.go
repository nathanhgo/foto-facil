package nodes

import (
	"errors"
	"image"
	"math/rand"
)

// AugmentationNode applies randomized (but reproducible, given a seed)
// transformations to each image in the batch — rotation, flips and mild
// noise — to synthetically expand a dataset for training more robust
// models. Each image gets its own derived seed so a batch produces varied
// but deterministic augmentations across runs.
type AugmentationNode struct {
	ID   string
	Seed int64 // Defaults to 42 when unset
}

func (n *AugmentationNode) GetID() string { return n.ID }

func (n *AugmentationNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("AugmentationNode requires at least one input image")
	}

	seed := n.Seed
	if seed == 0 {
		seed = 42
	}

	var processed []image.Image
	for i, img := range ctx.Images {
		rng := rand.New(rand.NewSource(seed + int64(i)))
		result := img

		angles := []RotateAngle{Rotate90, Rotate180, Rotate270}
		doRotate := rng.Intn(4) != 0 // 75% chance of rotating
		angle := angles[rng.Intn(len(angles))]
		doFlip := rng.Intn(2) == 0
		flipAxis := FlipAxis(rng.Intn(3))

		rotateCtx := &ProcessContext{Images: []image.Image{result}}
		rotateNode := &RotateFlipNode{
			ID:       n.ID + "_rotate",
			Angle:    angle,
			DoRotate: doRotate,
			DoFlip:   doFlip,
			Flip:     flipAxis,
		}
		if err := rotateNode.Process(rotateCtx); err != nil {
			return err
		}
		result = rotateCtx.Images[0]

		noiseCtx := &ProcessContext{Images: []image.Image{result}}
		noiseNode := &NoiseNode{ID: n.ID + "_noise", Type: "gaussian", Amount: 8, Seed: seed + int64(i) + 1}
		if err := noiseNode.Process(noiseCtx); err != nil {
			return err
		}
		result = noiseCtx.Images[0]

		processed = append(processed, result)
	}

	ctx.Images = processed
	return nil
}
