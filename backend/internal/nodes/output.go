package nodes

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
)

// OutputNode is responsible for saving processed images to disk
type OutputNode struct {
	ID        string
	OutputDir string
}

func (n *OutputNode) GetID() string {
	return n.ID
}

func (n *OutputNode) Process(ctx *ProcessContext) error {
	if len(ctx.Images) == 0 {
		return errors.New("OutputNode received no images to save")
	}

	if n.OutputDir == "" {
		n.OutputDir = "./output"
	}

	err := os.MkdirAll(n.OutputDir, os.ModePerm)
	if err != nil {
		return err
	}

	for i, img := range ctx.Images {
		filename := fmt.Sprintf("output_%s_%d.png", n.ID, i)
		path := filepath.Join(n.OutputDir, filename)

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
