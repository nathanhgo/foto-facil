package nodes

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	// Registra os decoders de JPEG, PNG e GIF automaticamente
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	// Suporte adicional a formatos
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// ProcessContext holds the data flowing between nodes
type ProcessContext struct {
	Images         []image.Image // Supports batch processing
	OriginalImages []image.Image // Copy of original images for comparison
	NodeOutputs    map[string][]image.Image // Map of node ID to its output images
}

// Node represents a process block in the DAG
type Node interface {
	GetID() string
	Process(ctx *ProcessContext) error
}

// InputNode reads real image files or directories from disk
type InputNode struct {
	ID        string
	FilePaths []string
}

func (n *InputNode) GetID() string { return n.ID }

func (n *InputNode) loadImage(path string, ctx *ProcessContext) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("InputNode: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("InputNode: cannot decode image %q: %w", path, err)
	}

	ctx.Images = append(ctx.Images, img)
	ctx.OriginalImages = append(ctx.OriginalImages, img)
	return nil
}

func (n *InputNode) Process(ctx *ProcessContext) error {
	if len(n.FilePaths) == 0 {
		return errors.New("no file paths provided to InputNode")
	}

	for _, path := range n.FilePaths {
		if path == "" {
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("InputNode: cannot stat path %q: %w", path, err)
		}

		if info.IsDir() {
			files, err := os.ReadDir(path)
			if err != nil {
				return fmt.Errorf("InputNode: cannot read directory %q: %w", path, err)
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				filePath := filepath.Join(path, file.Name())
				ext := strings.ToLower(filepath.Ext(filePath))
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" || ext == ".tiff" || ext == ".webp" {
					if err := n.loadImage(filePath, ctx); err != nil {
						log.Printf("InputNode: skipping invalid image in directory: %v", err)
					}
				}
			}
		} else {
			if err := n.loadImage(path, ctx); err != nil {
				return err
			}
		}
	}

	if len(ctx.Images) == 0 {
		return errors.New("InputNode: no valid images were loaded")
	}

	return nil
}
