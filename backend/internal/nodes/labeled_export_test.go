package nodes

import (
	"image"
	"os"
	"path/filepath"
	"testing"
)

func TestLabeledExportNodeSplitsIntoTrainValTest(t *testing.T) {
	tmpDir := t.TempDir()

	var images []image.Image
	for i := 0; i < 10; i++ {
		images = append(images, solidColorImage(2, 2, 10, 10, 10))
	}

	ctx := &ProcessContext{Images: images}
	node := &LabeledExportNode{ID: "export1", OutputDir: tmpDir, TrainRatio: 0.7, ValRatio: 0.2}

	if err := node.Process(ctx); err != nil {
		t.Fatalf("LabeledExportNode failed: %v", err)
	}

	trainFiles, _ := os.ReadDir(filepath.Join(tmpDir, "train"))
	valFiles, _ := os.ReadDir(filepath.Join(tmpDir, "val"))
	testFiles, _ := os.ReadDir(filepath.Join(tmpDir, "test"))

	if len(trainFiles) != 7 {
		t.Errorf("expected 7 train files, got %d", len(trainFiles))
	}
	if len(valFiles) != 2 {
		t.Errorf("expected 2 val files, got %d", len(valFiles))
	}
	if len(testFiles) != 1 {
		t.Errorf("expected 1 test file, got %d", len(testFiles))
	}
}

func TestLabeledExportNodeRequiresInputImage(t *testing.T) {
	ctx := &ProcessContext{}
	node := &LabeledExportNode{ID: "export2"}

	if err := node.Process(ctx); err == nil {
		t.Error("expected an error when no images are provided")
	}
}
