package nodes

import (
	"errors"
	"image"
	"image/color"
	"math"
)

// NormalizationNode rescales pixel intensities to prepare images for machine
// learning pipelines. "minmax" stretches each channel's observed range to
// fill 0-255; "zscore" centers each channel on its mean and maps +/-3
// standard deviations to 0-255. Real ML pipelines typically want float32
// tensors rather than 8-bit images — exporting that format directly (e.g.
// as .npy) is a follow-up; see mvp.md.
type NormalizationNode struct {
	ID     string
	Method string // "minmax" (default) or "zscore"
}

func (n *NormalizationNode) GetID() string { return n.ID }

// channelStats accumulates the min, max, mean and standard deviation of a
// single color channel across an image, then maps individual samples.
type channelStats struct {
	min, max   uint8
	sum, sumSq float64
	count      int
}

func newChannelStats() *channelStats {
	return &channelStats{min: 255, max: 0}
}

func (c *channelStats) observe(v uint8) {
	if v < c.min {
		c.min = v
	}
	if v > c.max {
		c.max = v
	}
	c.sum += float64(v)
	c.sumSq += float64(v) * float64(v)
	c.count++
}

func (c *channelStats) mean() float64 {
	if c.count == 0 {
		return 0
	}
	return c.sum / float64(c.count)
}

func (c *channelStats) stdDev() float64 {
	m := c.mean()
	v := c.sumSq/float64(imax(c.count, 1)) - m*m
	if v < 0 {
		v = 0
	}
	return math.Sqrt(v)
}

func (c *channelStats) minMaxByte(v uint8) uint8 {
	rangeV := float64(c.max) - float64(c.min)
	if rangeV == 0 {
		return v
	}
	return uint8(clampFloat((float64(v) - float64(c.min)) / rangeV * 255))
}

func (c *channelStats) zscoreByte(v uint8) uint8 {
	std := c.stdDev()
	if std == 0 {
		return 128
	}
	z := (float64(v) - c.mean()) / std
	return uint8(clampFloat((z + 3) / 6 * 255))
}

func (n *NormalizationNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("NormalizationNode requires at least one input image")
	}

	method := n.Method
	if method == "" {
		method = "minmax"
	}

	var processed []image.Image
	for _, img := range ctx.Images {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()

		var stats [3]*channelStats
		for i := range stats {
			stats[i] = newChannelStats()
		}

		type rgba struct{ r, g, bch, a uint8 }
		raw := make([]rgba, w*h)

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, bch, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
				px := rgba{uint8(r >> 8), uint8(g >> 8), uint8(bch >> 8), uint8(a >> 8)}
				raw[y*w+x] = px
				stats[0].observe(px.r)
				stats[1].observe(px.g)
				stats[2].observe(px.bch)
			}
		}

		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				px := raw[y*w+x]
				var out [3]uint8
				vals := [3]uint8{px.r, px.g, px.bch}
				for i, v := range vals {
					if method == "zscore" {
						out[i] = stats[i].zscoreByte(v)
					} else {
						out[i] = stats[i].minMaxByte(v)
					}
				}
				dst.Set(x, y, color.RGBA{out[0], out[1], out[2], px.a})
			}
		}
		processed = append(processed, dst)
	}

	ctx.Images = processed
	return nil
}
