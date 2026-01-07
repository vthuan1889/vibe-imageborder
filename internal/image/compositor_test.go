package image

import (
	"image"
	"image/color"
	"path/filepath"
	"testing"
)

func TestNewCompositor(t *testing.T) {
	svc := NewService()
	comp := NewCompositor(svc)
	if comp == nil {
		t.Error("NewCompositor returned nil")
	}
}

func TestComposite(t *testing.T) {
	svc := NewService()
	comp := NewCompositor(svc)

	// Create product image (white 100x100)
	product := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			product.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// Create frame image (200x200 with blue border, transparent center)
	frame := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for y := 0; y < 200; y++ {
		for x := 0; x < 200; x++ {
			if x < 20 || x > 180 || y < 20 || y > 180 {
				frame.Set(x, y, color.RGBA{0, 0, 255, 255}) // Blue border
			} else {
				frame.Set(x, y, color.RGBA{0, 0, 0, 0}) // Transparent center
			}
		}
	}

	result := comp.Composite(product, frame, "#f1eeea")

	if result.Width != 200 || result.Height != 200 {
		t.Errorf("Expected 200x200, got %dx%d", result.Width, result.Height)
	}

	if result.Image == nil {
		t.Error("Result image is nil")
	}
}

func TestCompositeWithPosition(t *testing.T) {
	svc := NewService()
	comp := NewCompositor(svc)

	product := image.NewRGBA(image.Rect(0, 0, 50, 50))
	frame := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Test "below" position
	result := comp.CompositeWithPosition(product, frame, "", "below")
	if result.Width != 100 || result.Height != 100 {
		t.Errorf("Expected 100x100, got %dx%d", result.Width, result.Height)
	}

	// Test default position
	result2 := comp.CompositeWithPosition(product, frame, "#ffffff", "above")
	if result2.Width != 100 || result2.Height != 100 {
		t.Errorf("Expected 100x100, got %dx%d", result2.Width, result2.Height)
	}
}

func TestCompositeNoBgColor(t *testing.T) {
	svc := NewService()
	comp := NewCompositor(svc)

	product := image.NewRGBA(image.Rect(0, 0, 50, 50))
	frame := image.NewRGBA(image.Rect(0, 0, 100, 100))

	result := comp.Composite(product, frame, "")
	if result.Width != 100 || result.Height != 100 {
		t.Errorf("Expected 100x100, got %dx%d", result.Width, result.Height)
	}
}

func TestToRGBA(t *testing.T) {
	// Test with already RGBA image
	rgba := image.NewRGBA(image.Rect(0, 0, 10, 10))
	result := ToRGBA(rgba)
	if result != rgba {
		t.Error("Expected same pointer for RGBA input")
	}

	// Test with non-RGBA image
	nrgba := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	result2 := ToRGBA(nrgba)
	if result2 == nil {
		t.Error("ToRGBA returned nil for NRGBA input")
	}
	if result2.Bounds() != nrgba.Bounds() {
		t.Error("ToRGBA changed bounds")
	}
}

func TestCompositeSaveResult(t *testing.T) {
	svc := NewService()
	comp := NewCompositor(svc)

	product := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			product.Set(x, y, color.RGBA{200, 200, 200, 255})
		}
	}

	frame := image.NewRGBA(image.Rect(0, 0, 150, 150))

	result := comp.Composite(product, frame, "#ffffff")

	// Save result to verify it's valid
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "composite.png")

	err := svc.SaveImage(result.Image, outputPath, "png", 90)
	if err != nil {
		t.Fatalf("Failed to save composite result: %v", err)
	}

	// Load and verify
	loaded, err := svc.LoadImage(outputPath)
	if err != nil {
		t.Fatalf("Failed to load composite result: %v", err)
	}

	w, h := svc.GetDimensions(loaded)
	if w != 150 || h != 150 {
		t.Errorf("Expected 150x150, got %dx%d", w, h)
	}
}
