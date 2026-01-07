package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestServiceLoadTemplate(t *testing.T) {
	content := `{
		"background": "#ffffff",
		"barcode": {
			"text": "[barcode]",
			"position": "90,1852",
			"fontsize": "50",
			"color": "white"
		}
	}`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	svc := NewService()

	// First load
	config1, err := svc.LoadTemplate(tmpFile)
	if err != nil {
		t.Fatalf("LoadTemplate failed: %v", err)
	}

	// Second load should return cached
	config2, err := svc.LoadTemplate(tmpFile)
	if err != nil {
		t.Fatalf("LoadTemplate (cached) failed: %v", err)
	}

	if config1 != config2 {
		t.Error("Expected same pointer from cache")
	}
}

func TestServiceGetFields(t *testing.T) {
	content := `{
		"barcode": {
			"text": "[barcode]",
			"position": "90,1852",
			"fontsize": "50",
			"color": "white"
		}
	}`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	svc := NewService()
	fields, err := svc.GetFields(tmpFile)
	if err != nil {
		t.Fatalf("GetFields failed: %v", err)
	}

	if len(fields) != 1 || fields[0] != "barcode" {
		t.Errorf("Expected [barcode], got %v", fields)
	}
}

func TestServiceGetOverlays(t *testing.T) {
	content := `{
		"barcode": {
			"text": "[barcode]",
			"position": "90,1852",
			"fontsize": "50",
			"color": "white"
		}
	}`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	svc := NewService()
	values := map[string]string{"barcode": "TEST123"}
	overlays, err := svc.GetOverlays(tmpFile, values)
	if err != nil {
		t.Fatalf("GetOverlays failed: %v", err)
	}

	if overlays["barcode"].Text != "TEST123" {
		t.Errorf("Expected TEST123, got %s", overlays["barcode"].Text)
	}
}

func TestServiceClearCache(t *testing.T) {
	svc := NewService()
	svc.ClearCache()

	if len(svc.cache) != 0 {
		t.Error("Expected empty cache after ClearCache")
	}
}
