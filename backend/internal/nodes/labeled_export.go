package nodes

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
)

// LabeledExportNode saves a batch of images split into train/val/test
// subfolders, the directory layout most ML training scripts expect. Splits
// are assigned deterministically by index using the given ratios; inferring
// labels from original filenames is a follow-up (see mvp.md) since the
// pipeline does not currently propagate source filenames through
// processing nodes.
type LabeledExportNode struct {
	ID         string
	OutputDir  string
	TrainRatio float64 // default 0.7
	ValRatio   float64 // default 0.15 (remainder goes to "test")
}

func (n *LabeledExportNode) GetID() string { return n.ID }

func (n *LabeledExportNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("LabeledExportNode received no images to save")
	}

	outputDir := n.OutputDir
	if outputDir == "" {
		outputDir = "./output"
	}
	trainRatio := n.TrainRatio
	if trainRatio <= 0 {
		trainRatio = 0.7
	}
	valRatio := n.ValRatio
	if valRatio <= 0 {
		valRatio = 0.15
	}

	total := len(ctx.Images)
	trainCount := int(float64(total) * trainRatio)
	valCount := int(float64(total) * valRatio)

	for i, img := range ctx.Images {
		split := "test"
		if i < trainCount {
			split = "train"
		} else if i < trainCount+valCount {
			split = "val"
		}

		dir := filepath.Join(outputDir, split)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		filename := fmt.Sprintf("%s_%d.png", n.ID, i)
		path := filepath.Join(dir, filename)
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		err = png.Encode(file, img)
		file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
