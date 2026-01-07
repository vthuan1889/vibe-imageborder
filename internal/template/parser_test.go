package template

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestParseTemplate(t *testing.T) {
	content := `{
		"background": "#f1eeea",
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

	config, err := ParseTemplate(tmpFile)
	if err != nil {
		t.Fatalf("ParseTemplate failed: %v", err)
	}

	if config.Background != "#f1eeea" {
		t.Errorf("Expected background #f1eeea, got %s", config.Background)
	}

	if len(config.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(config.Fields))
	}

	if config.Fields["barcode"].Text != "[barcode]" {
		t.Errorf("Expected [barcode], got %s", config.Fields["barcode"].Text)
	}
}

func TestExtractFields(t *testing.T) {
	content := `{
		"barcode": {
			"text": "[barcode]",
			"position": "90,1852",
			"fontsize": "50",
			"color": "white"
		},
		"size": {
			"text": "D[size_dai] x R[size_rong] x C[size_cao] CM",
			"position": "1100,10",
			"fontsize": "60",
			"color": "white"
		}
	}`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	config, err := ParseTemplate(tmpFile)
	if err != nil {
		t.Fatalf("ParseTemplate failed: %v", err)
	}

	fields := ExtractFields(config)
	sort.Strings(fields)

	expected := []string{"barcode", "size_cao", "size_dai", "size_rong"}
	sort.Strings(expected)

	if len(fields) != len(expected) {
		t.Errorf("Expected %d fields, got %d: %v", len(expected), len(fields), fields)
	}

	for i, f := range fields {
		if f != expected[i] {
			t.Errorf("Expected field %s, got %s", expected[i], f)
		}
	}
}

func TestApplyValues(t *testing.T) {
	content := `{
		"barcode": {
			"text": "[barcode]",
			"position": "90,1852",
			"fontsize": "50",
			"color": "white"
		},
		"price": {
			"text": "Giá [price]K",
			"position": "10,1712",
			"fontsize": "50",
			"color": "white"
		}
	}`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	config, err := ParseTemplate(tmpFile)
	if err != nil {
		t.Fatalf("ParseTemplate failed: %v", err)
	}

	values := map[string]string{
		"barcode": "ABC123",
		"price":   "500",
	}
	result := ApplyValues(config, values)

	if result["barcode"].Text != "ABC123" {
		t.Errorf("Expected ABC123, got %s", result["barcode"].Text)
	}

	if result["price"].Text != "Giá 500K" {
		t.Errorf("Expected 'Giá 500K', got %s", result["price"].Text)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	content := `{ invalid json }`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.txt")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err := ParseTemplate(tmpFile)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestParseNonexistentFile(t *testing.T) {
	_, err := ParseTemplate("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}
