package image

import (
	"image"
	"image/color"
	"path/filepath"
	"testing"
)

func TestNewService(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Error("NewService returned nil")
	}
}

func TestSaveAndLoadImage(t *testing.T) {
	svc := NewService()

	// Create test image
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			testImg.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test.png")

	// Save
	err := svc.SaveImage(testImg, testPath, "png", 90)
	if err != nil {
		t.Fatalf("SaveImage failed: %v", err)
	}

	// Load
	loaded, err := svc.LoadImage(testPath)
	if err != nil {
		t.Fatalf("LoadImage failed: %v", err)
	}

	w, h := svc.GetDimensions(loaded)
	if w != 100 || h != 100 {
		t.Errorf("Expected 100x100, got %dx%d", w, h)
	}
}

func TestSaveImageJPG(t *testing.T) {
	svc := NewService()

	testImg := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			testImg.Set(x, y, color.RGBA{0, 255, 0, 255})
		}
	}

	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test.jpg")

	err := svc.SaveImage(testImg, testPath, "jpg", 85)
	if err != nil {
		t.Fatalf("SaveImage (JPG) failed: %v", err)
	}

	loaded, err := svc.LoadImage(testPath)
	if err != nil {
		t.Fatalf("LoadImage (JPG) failed: %v", err)
	}

	w, h := svc.GetDimensions(loaded)
	if w != 50 || h != 50 {
		t.Errorf("Expected 50x50, got %dx%d", w, h)
	}
}

func TestResizeToFit(t *testing.T) {
	svc := NewService()

	// Create 200x100 image (2:1 aspect)
	testImg := image.NewRGBA(image.Rect(0, 0, 200, 100))

	// Resize to fit 100x100 (should become 100x50)
	resized := svc.ResizeToFit(testImg, 100, 100)
	w, h := svc.GetDimensions(resized)

	if w != 100 || h != 50 {
		t.Errorf("Expected 100x50, got %dx%d", w, h)
	}
}

func TestResizeToFill(t *testing.T) {
	svc := NewService()

	// Create 200x100 image
	testImg := image.NewRGBA(image.Rect(0, 0, 200, 100))

	// Resize to fill 100x100 (should crop to 100x100)
	resized := svc.ResizeToFill(testImg, 100, 100)
	w, h := svc.GetDimensions(resized)

	if w != 100 || h != 100 {
		t.Errorf("Expected 100x100, got %dx%d", w, h)
	}
}

func TestCreateBlankCanvas(t *testing.T) {
	svc := NewService()

	canvas := svc.CreateBlankCanvas(150, 150, "#ff0000")
	w, h := svc.GetDimensions(canvas)

	if w != 150 || h != 150 {
		t.Errorf("Expected 150x150, got %dx%d", w, h)
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		hex      string
		expected color.RGBA
	}{
		{"#ff0000", color.RGBA{255, 0, 0, 255}},
		{"ff0000", color.RGBA{255, 0, 0, 255}},
		{"#00ff00", color.RGBA{0, 255, 0, 255}},
		{"#0000ff", color.RGBA{0, 0, 255, 255}},
		{"#ffffff", color.RGBA{255, 255, 255, 255}},
		{"invalid", color.RGBA{255, 255, 255, 255}}, // fallback to white
	}

	for _, tc := range tests {
		result := parseColor(tc.hex)
		if r, ok := result.(color.RGBA); ok {
			if r != tc.expected {
				t.Errorf("parseColor(%s) = %v, expected %v", tc.hex, r, tc.expected)
			}
		}
	}
}

func TestUnsupportedFormat(t *testing.T) {
	svc := NewService()
	testImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test.xyz")

	err := svc.SaveImage(testImg, testPath, "xyz", 90)
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
}
