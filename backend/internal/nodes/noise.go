package nodes

import (
	"errors"
	"image"
	"image/color"
	"math/rand"
)

// NoiseNode injects synthetic sensor-like noise into images, useful for
// building datasets that test the robustness of downstream algorithms
// (e.g. simulating a noisy sensor to validate the StackingNode).
type NoiseNode struct {
	ID     string
	Type   string  // "gaussian" (default), "salt_pepper" or "speckle"
	Amount float64 // Std-dev (0-255 scale) for gaussian/speckle, or probability (0-1) for salt_pepper
	Seed   int64   // Deterministic seed; defaults to 1 when unset so runs are reproducible
}

func (n *NoiseNode) GetID() string { return n.ID }

func (n *NoiseNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("NoiseNode requires at least one input image")
	}

	noiseType := n.Type
	if noiseType == "" {
		noiseType = "gaussian"
	}
	amount := n.Amount
	if amount <= 0 {
		if noiseType == "salt_pepper" {
			amount = 0.02
		} else {
			amount = 20
		}
	}

	seed := n.Seed
	if seed == 0 {
		seed = 1
	}
	rng := rand.New(rand.NewSource(seed))

	var processed []image.Image
	for _, img := range ctx.Images {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()
		dst := image.NewRGBA(image.Rect(0, 0, w, h))

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, bch, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
				r8, g8, b8, a8 := float64(r>>8), float64(g>>8), float64(bch>>8), uint8(a>>8)

				switch noiseType {
				case "salt_pepper":
					roll := rng.Float64()
					if roll < amount/2 {
						r8, g8, b8 = 0, 0, 0
					} else if roll < amount {
						r8, g8, b8 = 255, 255, 255
					}
				case "speckle":
					factor := 1 + (rng.Float64()*2-1)*(amount/255)
					r8 *= factor
					g8 *= factor
					b8 *= factor
				default: // gaussian
					r8 += rng.NormFloat64() * amount
					g8 += rng.NormFloat64() * amount
					b8 += rng.NormFloat64() * amount
				}

				dst.Set(x, y, color.RGBA{uint8(clampFloat(r8)), uint8(clampFloat(g8)), uint8(clampFloat(b8)), a8})
			}
		}
		processed = append(processed, dst)
	}

	ctx.Images = processed
	return nil
}
