package image

import (
	"image/color"
	"os"
	"testing"

	"github.com/disintegration/imaging"
)

func TestLoadImage(t *testing.T) {
	svc := NewService()

	// Create output dir if not exists
	os.MkdirAll("../../tests/output", 0755)

	// Create test image
	tmpImg := imaging.New(800, 600, color.RGBA{255, 0, 0, 255})
	tmpPath := "../../tests/output/test_load.png"
	imaging.Save(tmpImg, tmpPath)

	// Test load
	img, err := svc.LoadImage(tmpPath)
	if err != nil {
		t.Fatalf("LoadImage failed: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 800 || bounds.Dy() != 600 {
		t.Errorf("Expected 800x600, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestLoadImageTooLarge(t *testing.T) {
	svc := NewService()

	// Create output dir if not exists
	os.MkdirAll("../../tests/output", 0755)

	// Create very large test image (should fail)
	largeImg := imaging.New(11000, 11000, color.RGBA{0, 255, 0, 255})
	largePath := "../../tests/output/test_large.png"
	imaging.Save(largeImg, largePath)

	// Test load (should fail)
	_, err := svc.LoadImage(largePath)
	if err == nil {
		t.Error("Expected error for oversized image, got nil")
	}

	// Cleanup
	os.Remove(largePath)
}

func TestResizeToFit(t *testing.T) {
	svc := NewService()

	// Create test image: 1000x500 (2:1 aspect ratio)
	img := imaging.New(1000, 500, color.RGBA{0, 255, 0, 255})

	// Fit into 800x800 (should be 800x400 to preserve 2:1)
	resized := svc.ResizeToFit(img, 800, 800)

	bounds := resized.Bounds()
	expectedW, expectedH := 800, 400

	if bounds.Dx() != expectedW || bounds.Dy() != expectedH {
		t.Errorf("Expected %dx%d, got %dx%d",
			expectedW, expectedH, bounds.Dx(), bounds.Dy())
	}
}

func TestCompositeImages(t *testing.T) {
	svc := NewService()

	// Create output dir if not exists
	os.MkdirAll("../../tests/output", 0755)

	// Create test images
	frame := imaging.New(1000, 1000, color.RGBA{200, 200, 200, 255})
	product := imaging.New(600, 400, color.RGBA{255, 0, 0, 255})

	framePath := "../../tests/output/test_frame.png"
	productPath := "../../tests/output/test_product.png"

	imaging.Save(frame, framePath)
	imaging.Save(product, productPath)

	// Test composite
	result, err := svc.CompositeImages(productPath, framePath)
	if err != nil {
		t.Fatalf("CompositeImages failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Composite unsuccessful: %s", result.ErrorMessage)
	}

	// Verify output dimensions match frame
	bounds := result.Image.Bounds()
	if bounds.Dx() != 1000 || bounds.Dy() != 1000 {
		t.Errorf("Expected 1000x1000, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Save composite for visual inspection
	svc.SaveImage(result.Image, "../../tests/output/test_composite_result.png")
}

func TestSaveImage(t *testing.T) {
	svc := NewService()

	// Create output dir if not exists
	os.MkdirAll("../../tests/output", 0755)

	// Create test image
	img := imaging.New(400, 300, color.RGBA{0, 0, 255, 255})
	path := "../../tests/output/test_save.png"

	// Test save
	err := svc.SaveImage(img, path)
	if err != nil {
		t.Fatalf("SaveImage failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Verify can load it back
	loaded, err := svc.LoadImage(path)
	if err != nil {
		t.Fatalf("Failed to reload saved image: %v", err)
	}

	bounds := loaded.Bounds()
	if bounds.Dx() != 400 || bounds.Dy() != 300 {
		t.Errorf("Reloaded image dimensions incorrect: got %dx%d", bounds.Dx(), bounds.Dy())
	}
}
