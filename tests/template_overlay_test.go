package tests

import (
	"os"
	"testing"

	"vibe-imageborder/internal/template"
)

func TestRealTemplateOverlay(t *testing.T) {
	// Test với template trong fixtures
	templatePath := "fixtures/templates/test-template.txt"

	// Skip if template doesn't exist
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Skip("Template fixture not found, skipping test")
	}

	templateSvc := template.NewService()

	// Load fields
	fields, err := templateSvc.GetFields(templatePath)
	if err != nil {
		t.Fatalf("Failed to get fields: %v", err)
	}

	t.Logf("Found fields: %v", fields)

	// User values
	values := map[string]string{
		"barcode":   "SP12345",
		"size_dai":  "100",
		"size_rong": "50",
	}

	// Get overlays with values applied
	overlays, err := templateSvc.GetOverlays(templatePath, values)
	if err != nil {
		t.Fatalf("Failed to get overlays: %v", err)
	}

	// Verify each overlay has correct data
	for key, overlay := range overlays {
		t.Logf("Overlay %s:", key)
		t.Logf("  Text: %s", overlay.Text)
		t.Logf("  Position: %s", overlay.Position)
		t.Logf("  FontSize: %d", overlay.FontSize)
		t.Logf("  Color: %s", overlay.Color)
	}

	// Check barcode overlay
	if barcode, ok := overlays["barcode"]; ok {
		if barcode.Text != "SP12345" {
			t.Errorf("Expected barcode text 'SP12345', got '%s'", barcode.Text)
		}
		if barcode.Position != "100,50" {
			t.Errorf("Expected position '100,50', got '%s'", barcode.Position)
		}
		if barcode.FontSize != 24 {
			t.Errorf("Expected fontsize 24, got %d", barcode.FontSize)
		}
		if barcode.Color != "black" {
			t.Errorf("Expected color 'black', got '%s'", barcode.Color)
		}
	} else {
		t.Error("barcode overlay not found")
	}

	// Check size_dai overlay
	if sizeDai, ok := overlays["size_dai"]; ok {
		if sizeDai.Text != "100" {
			t.Errorf("Expected size_dai text '100', got '%s'", sizeDai.Text)
		}
	}

	// Check size_rong overlay
	if sizeRong, ok := overlays["size_rong"]; ok {
		if sizeRong.Text != "50" {
			t.Errorf("Expected size_rong text '50', got '%s'", sizeRong.Text)
		}
	}
}
