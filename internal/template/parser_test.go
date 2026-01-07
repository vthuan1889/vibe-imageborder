package template

import (
	"os"
	"path/filepath"
	"testing"
	"vibe-imageborder/internal/models"
)

func TestParseTemplate(t *testing.T) {
	// Create test template
	testJSON := `{
		"barcode": {
			"text": "[barcode]",
			"position": "98,1720",
			"fontsize": "45",
			"color": "white"
		},
		"size": {
			"text": "D[size_dai] x R[size_rong] cm",
			"position": "26,1852",
			"fontsize": "40",
			"color": "#FFFFFF"
		}
	}`

	tmpFile := filepath.Join(t.TempDir(), "test_template.txt")
	os.WriteFile(tmpFile, []byte(testJSON), 0644)

	tmpl, err := ParseTemplate(tmpFile)
	if err != nil {
		t.Fatalf("ParseTemplate failed: %v", err)
	}

	if len(tmpl) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(tmpl))
	}

	if tmpl["barcode"].Text != "[barcode]" {
		t.Errorf("Unexpected barcode text: %s", tmpl["barcode"].Text)
	}
}

func TestExtractDynamicFields(t *testing.T) {
	tmpl := models.Template{
		"barcode": models.TemplateField{
			Text: "[barcode]",
		},
		"size": models.TemplateField{
			Text: "D[size_dai] x R[size_rong] x C[size_cao] cm",
		},
	}

	fields := ExtractDynamicFields(tmpl)

	expectedCount := 4 // barcode, size_dai, size_rong, size_cao
	if len(fields) != expectedCount {
		t.Errorf("Expected %d fields, got %d: %v", expectedCount, len(fields), fields)
	}
}

func TestReplaceVariables(t *testing.T) {
	text := "D[size_dai] x R[size_rong] cm"
	values := map[string]string{
		"size_dai":  "30",
		"size_rong": "20",
	}

	result := ReplaceVariables(text, values)
	expected := "D30 x R20 cm"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestRealTemplates(t *testing.T) {
	templates := []string{
		"../../tests/fixtures/templates/khung-002-05.txt",
		"../../tests/fixtures/templates/khung-004-01.txt",
	}

	for _, path := range templates {
		t.Run(filepath.Base(path), func(t *testing.T) {
			tmpl, err := ParseTemplate(path)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", path, err)
			}

			fields := ExtractDynamicFields(tmpl)
			if len(fields) == 0 {
				t.Errorf("No dynamic fields found in %s", path)
			}

			t.Logf("Template %s has %d fields: %v", path, len(fields), fields)
		})
	}
}

// Test error cases
func TestParseTemplateErrors(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "missing file",
			content:     "",
			expectError: true,
			errorMsg:    "failed to read template file",
		},
		{
			name:        "invalid JSON",
			content:     "not json",
			expectError: true,
			errorMsg:    "failed to parse template JSON",
		},
		{
			name: "empty template",
			content: `{}`,
			expectError: true,
			errorMsg:    "template is empty",
		},
		{
			name: "missing text field",
			content: `{
				"field1": {
					"position": "10,20",
					"fontsize": "40",
					"color": "white"
				}
			}`,
			expectError: true,
			errorMsg:    "text is empty",
		},
		{
			name: "missing position field",
			content: `{
				"field1": {
					"text": "test",
					"fontsize": "40",
					"color": "white"
				}
			}`,
			expectError: true,
			errorMsg:    "position is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "missing file" {
				_, err := ParseTemplate("/nonexistent/file.txt")
				if err == nil {
					t.Error("Expected error for missing file")
				}
				return
			}

			tmpFile := filepath.Join(t.TempDir(), "test.txt")
			os.WriteFile(tmpFile, []byte(tt.content), 0644)

			_, err := ParseTemplate(tmpFile)
			if tt.expectError && err == nil {
				t.Errorf("Expected error containing %q, got nil", tt.errorMsg)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Test ReplaceVariables edge cases
func TestReplaceVariablesEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		values   map[string]string
		expected string
	}{
		{
			name:     "no placeholders",
			text:     "Plain text",
			values:   map[string]string{},
			expected: "Plain text",
		},
		{
			name:     "missing value keeps placeholder",
			text:     "Price: [price]",
			values:   map[string]string{},
			expected: "Price: [price]",
		},
		{
			name:     "multiple same placeholder",
			text:     "[x] + [x] = 2*[x]",
			values:   map[string]string{"x": "5"},
			expected: "5 + 5 = 2*5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceVariables(tt.text, tt.values)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test Service methods
func TestService(t *testing.T) {
	svc := NewService()

	// Create test template
	testJSON := `{
		"field1": {
			"text": "[var1] and [var2]",
			"position": "10,20",
			"fontsize": "40",
			"color": "white"
		}
	}`

	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte(testJSON), 0644)

	// Test Load
	tmpl, err := svc.Load(tmpFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Test GetDynamicFields
	fields := svc.GetDynamicFields(tmpl)
	if len(fields) != 2 {
		t.Errorf("Expected 2 fields, got %d: %v", len(fields), fields)
	}

	// Test ApplyValues
	values := map[string]string{
		"var1": "Hello",
		"var2": "World",
	}
	result := svc.ApplyValues(tmpl, values)
	if result["field1"].Text != "Hello and World" {
		t.Errorf("ApplyValues failed: %s", result["field1"].Text)
	}
}
