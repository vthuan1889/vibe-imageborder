package tests

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	imgservice "vibe-imageborder/internal/image"
	"vibe-imageborder/internal/template"
)

func createTestImage(width, height int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func saveTestImage(t *testing.T, img image.Image, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}
}

func TestIntegration_BasicComposite(t *testing.T) {
	// Setup test images
	productPath := filepath.Join("fixtures", "products", "test-product.png")
	framePath := filepath.Join("fixtures", "frames", "test-frame.png")
	outputPath := filepath.Join("output", "test-composite")

	// Create test product (400x400 red square)
	product := createTestImage(400, 400, color.RGBA{255, 0, 0, 255})
	saveTestImage(t, product, productPath)

	// Create test frame (500x500 with transparent center)
	frame := image.NewRGBA(image.Rect(0, 0, 500, 500))
	for y := 0; y < 500; y++ {
		for x := 0; x < 500; x++ {
			// Border area = blue, center = transparent
			if x < 50 || x >= 450 || y < 50 || y >= 450 {
				frame.Set(x, y, color.RGBA{0, 0, 255, 255})
			} else {
				frame.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
	saveTestImage(t, frame, framePath)

	// Test
	imageSvc := imgservice.NewService()
	compositor := imgservice.NewCompositor(imageSvc)

	loadedProduct, err := imageSvc.LoadImage(productPath)
	if err != nil {
		t.Fatalf("Failed to load product: %v", err)
	}

	loadedFrame, err := imageSvc.LoadImage(framePath)
	if err != nil {
		t.Fatalf("Failed to load frame: %v", err)
	}

	result := compositor.Composite(loadedProduct, loadedFrame, "#f1eeea")

	if result.Width != 500 || result.Height != 500 {
		t.Errorf("Expected 500x500, got %dx%d", result.Width, result.Height)
	}

	err = imageSvc.SaveImage(result.Image, outputPath, "png", 90)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath + ".png"); os.IsNotExist(err) {
		t.Fatalf("Output file not created at %s.png", outputPath)
	}
}

func TestIntegration_TemplateLoading(t *testing.T) {
	templatePath := filepath.Join("fixtures", "templates", "test-template.txt")

	templateSvc := template.NewService()

	// Load fields
	fields, err := templateSvc.GetFields(templatePath)
	if err != nil {
		t.Fatalf("Failed to get fields: %v", err)
	}

	expectedFields := []string{"barcode", "size_dai", "size_rong"}
	if len(fields) != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), len(fields))
	}

	// Test overlay generation
	values := map[string]string{
		"barcode":   "TEST001",
		"size_dai":  "100",
		"size_rong": "50",
	}

	overlays, err := templateSvc.GetOverlays(templatePath, values)
	if err != nil {
		t.Fatalf("Failed to get overlays: %v", err)
	}

	if len(overlays) != 3 {
		t.Errorf("Expected 3 overlays, got %d", len(overlays))
	}

	// Verify values are applied
	if overlay, ok := overlays["barcode"]; ok {
		if overlay.Text != "TEST001" {
			t.Errorf("Expected barcode text 'TEST001', got '%s'", overlay.Text)
		}
	} else {
		t.Error("barcode overlay not found")
	}
}

func TestIntegration_BackgroundColor(t *testing.T) {
	templatePath := filepath.Join("fixtures", "templates", "test-template.txt")

	templateSvc := template.NewService()

	bg, err := templateSvc.GetBackground(templatePath)
	if err != nil {
		t.Fatalf("Failed to get background: %v", err)
	}

	if bg != "#f1eeea" {
		t.Errorf("Expected background '#f1eeea', got '%s'", bg)
	}
}

func TestIntegration_OutputFormats(t *testing.T) {
	productPath := filepath.Join("fixtures", "products", "test-product.png")
	framePath := filepath.Join("fixtures", "frames", "test-frame.png")

	imageSvc := imgservice.NewService()
	compositor := imgservice.NewCompositor(imageSvc)

	product, _ := imageSvc.LoadImage(productPath)
	frame, _ := imageSvc.LoadImage(framePath)

	result := compositor.Composite(product, frame, "")

	formats := []struct {
		format  string
		quality int
	}{
		{"png", 0},
		{"jpg", 90},
		{"jpg", 50},
	}

	for _, f := range formats {
		outputPath := filepath.Join("output", "test-format-"+f.format)
		err := imageSvc.SaveImage(result.Image, outputPath, f.format, f.quality)
		if err != nil {
			t.Errorf("Failed to save %s: %v", f.format, err)
		}
	}
}
