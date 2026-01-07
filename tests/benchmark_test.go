package tests

import (
	"os"
	"path/filepath"
	"testing"

	"vibe-imageborder/internal/image"
	"vibe-imageborder/internal/template"
)

// BenchmarkSingleImage benchmarks single image processing
func BenchmarkSingleImage(b *testing.B) {
	imgSvc := image.NewService()
	tmplSvc := template.NewService()

	// Load template
	tmplPath := filepath.Join("fixtures", "templates", "khung-002-05.txt")
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		b.Skip("Template not found for benchmark")
		return
	}

	tmpl, err := tmplSvc.Load(tmplPath)
	if err != nil {
		b.Fatalf("Failed to load template: %v", err)
	}

	fieldValues := map[string]string{
		"barcode":   "BENCH123",
		"size_dai":  "30",
		"size_rong": "20",
		"size_cao":  "15",
	}

	productPath := filepath.Join("fixtures", "products", "product-01.jpg")
	framePath := filepath.Join("fixtures", "frames", "frame-01.png")

	if _, err := os.Stat(productPath); os.IsNotExist(err) {
		b.Skip("Product image not found for benchmark")
		return
	}
	if _, err := os.Stat(framePath); os.IsNotExist(err) {
		b.Skip("Frame image not found for benchmark")
		return
	}

	outputDir := "output"
	os.MkdirAll(outputDir, 0755)
	outputPath := filepath.Join(outputDir, "bench_output.png")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := imgSvc.ProcessSingle(
			productPath,
			framePath,
			outputPath,
			tmpl,
			fieldValues,
		)
		if err != nil {
			b.Fatalf("ProcessSingle failed: %v", err)
		}
	}
}

// BenchmarkTemplateLoading benchmarks template loading and parsing
func BenchmarkTemplateLoading(b *testing.B) {
	tmplSvc := template.NewService()
	tmplPath := filepath.Join("fixtures", "templates", "khung-002-05.txt")

	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		b.Skip("Template not found for benchmark")
		return
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tmplSvc.Load(tmplPath)
		if err != nil {
			b.Fatalf("Template load failed: %v", err)
		}
	}
}

// BenchmarkImageLoading benchmarks image loading
func BenchmarkImageLoading(b *testing.B) {
	imgSvc := image.NewService()
	productPath := filepath.Join("fixtures", "products", "product-01.jpg")

	if _, err := os.Stat(productPath); os.IsNotExist(err) {
		b.Skip("Product image not found for benchmark")
		return
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := imgSvc.LoadImage(productPath)
		if err != nil {
			b.Fatalf("Image load failed: %v", err)
		}
	}
}
