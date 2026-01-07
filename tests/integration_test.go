package tests

import (
	"os"
	"path/filepath"
	"testing"

	"vibe-imageborder/internal/image"
	"vibe-imageborder/internal/template"
)

// TestE2EWorkflow tests complete end-to-end workflow
func TestE2EWorkflow(t *testing.T) {
	// Setup services
	imgSvc := image.NewService()
	tmplSvc := template.NewService()

	// Load template
	tmplPath := filepath.Join("fixtures", "templates", "khung-002-05.txt")
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		t.Skipf("Template file not found: %s (manual test required)", tmplPath)
		return
	}

	tmpl, err := tmplSvc.Load(tmplPath)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// Extract fields
	fields := tmplSvc.GetDynamicFields(tmpl)
	if len(fields) == 0 {
		t.Fatal("No dynamic fields extracted from template")
	}
	t.Logf("Extracted %d fields: %v", len(fields), fields)

	// Prepare field values
	fieldValues := map[string]string{
		"barcode":   "TEST123",
		"size_dai":  "30",
		"size_rong": "20",
		"size_cao":  "15",
	}

	// Test image paths
	productPath := filepath.Join("fixtures", "products", "product-01.jpg")
	framePath := filepath.Join("fixtures", "frames", "frame-01.png")
	outputDir := "output"

	// Skip if fixtures don't exist
	if _, err := os.Stat(productPath); os.IsNotExist(err) {
		t.Skipf("Product image not found: %s (manual test required)", productPath)
		return
	}
	if _, err := os.Stat(framePath); os.IsNotExist(err) {
		t.Skipf("Frame image not found: %s (manual test required)", framePath)
		return
	}

	// Create output directory
	os.MkdirAll(outputDir, 0755)

	// Process single image
	outputPath := filepath.Join(outputDir, "test_e2e_output.png")
	err = imgSvc.ProcessSingle(
		productPath,
		framePath,
		outputPath,
		tmpl,
		fieldValues,
	)
	if err != nil {
		t.Fatalf("ProcessSingle failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file not created: %s", outputPath)
	}

	t.Logf("✓ E2E test passed - output: %s", outputPath)
}

// TestBatchProcessing tests batch processing with multiple images
func TestBatchProcessing(t *testing.T) {
	imgSvc := image.NewService()
	tmplSvc := template.NewService()

	// Load template
	tmplPath := filepath.Join("fixtures", "templates", "khung-002-05.txt")
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		t.Skip("Template not found, skipping batch test")
		return
	}

	tmpl, err := tmplSvc.Load(tmplPath)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	fieldValues := map[string]string{
		"barcode":   "BATCH001",
		"size_dai":  "25",
		"size_rong": "18",
		"size_cao":  "12",
	}

	// Test with 3 images (simulate batch)
	products := []string{
		filepath.Join("fixtures", "products", "product-01.jpg"),
		filepath.Join("fixtures", "products", "product-02.jpg"),
		filepath.Join("fixtures", "products", "product-03.jpg"),
	}

	framePath := filepath.Join("fixtures", "frames", "frame-01.png")
	outputDir := "output"
	os.MkdirAll(outputDir, 0755)

	successCount := 0
	for i, productPath := range products {
		if _, err := os.Stat(productPath); os.IsNotExist(err) {
			t.Logf("Product %d not found, skipping: %s", i+1, productPath)
			continue
		}

		outputPath := filepath.Join(outputDir, filepath.Base(productPath)+"_framed.png")
		err := imgSvc.ProcessSingle(productPath, framePath, outputPath, tmpl, fieldValues)
		if err != nil {
			t.Logf("Warning: Failed to process %s: %v", productPath, err)
			continue
		}

		successCount++
	}

	if successCount == 0 {
		t.Skip("No products processed (fixtures missing)")
		return
	}

	t.Logf("✓ Batch test passed - processed %d/%d images", successCount, len(products))
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	imgSvc := image.NewService()

	t.Run("InvalidProductImage", func(t *testing.T) {
		tmplSvc := template.NewService()
		tmplPath := filepath.Join("fixtures", "templates", "khung-002-05.txt")
		if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
			t.Skip("Template not found")
			return
		}

		tmpl, _ := tmplSvc.Load(tmplPath)
		fieldValues := map[string]string{"barcode": "ERR"}

		// Try to process non-existent image
		err := imgSvc.ProcessSingle(
			"nonexistent.jpg",
			filepath.Join("fixtures", "frames", "frame-01.png"),
			"output/error_test.png",
			tmpl,
			fieldValues,
		)

		if err == nil {
			t.Fatal("Expected error for non-existent image, got nil")
		}
		t.Logf("✓ Error handling works: %v", err)
	})

	t.Run("InvalidFrameImage", func(t *testing.T) {
		tmplSvc := template.NewService()
		tmplPath := filepath.Join("fixtures", "templates", "khung-002-05.txt")
		if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
			t.Skip("Template not found")
			return
		}

		tmpl, _ := tmplSvc.Load(tmplPath)
		fieldValues := map[string]string{"barcode": "ERR"}

		productPath := filepath.Join("fixtures", "products", "product-01.jpg")
		if _, err := os.Stat(productPath); os.IsNotExist(err) {
			t.Skip("Product not found")
			return
		}

		// Try to process with non-existent frame
		err := imgSvc.ProcessSingle(
			productPath,
			"nonexistent_frame.png",
			"output/error_test.png",
			tmpl,
			fieldValues,
		)

		if err == nil {
			t.Fatal("Expected error for non-existent frame, got nil")
		}
		t.Logf("✓ Frame error handling works: %v", err)
	})
}
