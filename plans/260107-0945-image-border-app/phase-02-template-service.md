# Phase 2: Template Service

**Goal:** Parse template JSON files, extract dynamic fields, replace variables

**Status:** ✅ COMPLETED (2026-01-07 12:15 PM)

**Completion Summary:** All tasks for Phase 2, including type definition, template parser implementation, template service creation, and unit testing, have been successfully completed and verified. The service is ready for integration with Phase 3.

**Duration:** ~2-3 hours

**Dependencies:** Phase 1 complete

**Review:** [Code Review Report](../reports/code-reviewer-260107-1213-phase2-template-service.md)

---

## Overview

Implement Template Service để:
1. Load và parse JSON template files
2. Extract dynamic field names `[barcode]`, `[size_dai]`, etc.
3. Replace placeholders với user input values

---

## Task 2.1: Define Types

Create `internal/models/types.go`:

```go
package models

// TemplateField represents a single text field trong template
type TemplateField struct {
	Text     string `json:"text"`     // Text với [placeholders]
	Position string `json:"position"` // "x,y" format
	FontSize string `json:"fontsize"` // "45"
	Color    string `json:"color"`    // "white", "#FFFFFF"
}

// Template is a map of field name → field config
type Template map[string]TemplateField

// CompositeRequest chứa tất cả data cho batch processing
type CompositeRequest struct {
	ProductPaths []string          `json:"productPaths"` // Paths to product images
	FramePath    string            `json:"framePath"`    // Path to frame image
	TemplatePath string            `json:"templatePath"` // Path to template JSON
	FieldValues  map[string]string `json:"fieldValues"`  // User input: "barcode" → "ABC123"
	OutputDir    string            `json:"outputDir"`    // Output directory
}
```

---

## Task 2.2: Implement Template Parser

Create `internal/template/parser.go`:

```go
package template

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"vibe-imageborder/internal/models"
)

// ParseTemplate loads và parses template JSON file
func ParseTemplate(path string) (models.Template, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Parse JSON
	var tmpl models.Template
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse template JSON: %w", err)
	}

	// Validate template
	if err := validateTemplate(tmpl); err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}

	return tmpl, nil
}

// validateTemplate checks template structure
func validateTemplate(tmpl models.Template) error {
	if len(tmpl) == 0 {
		return fmt.Errorf("template is empty")
	}

	for name, field := range tmpl {
		if field.Text == "" {
			return fmt.Errorf("field %s: text is empty", name)
		}
		if field.Position == "" {
			return fmt.Errorf("field %s: position is empty", name)
		}
		if field.FontSize == "" {
			return fmt.Errorf("field %s: fontsize is empty", name)
		}
		if field.Color == "" {
			return fmt.Errorf("field %s: color is empty", name)
		}
	}

	return nil
}

// ExtractDynamicFields finds all unique [field] placeholders trong template
func ExtractDynamicFields(tmpl models.Template) []string {
	fieldSet := make(map[string]bool)
	re := regexp.MustCompile(`\[([^\]]+)\]`)

	for _, field := range tmpl {
		matches := re.FindAllStringSubmatch(field.Text, -1)
		for _, match := range matches {
			if len(match) > 1 {
				fieldSet[match[1]] = true
			}
		}
	}

	// Convert set to slice
	fields := make([]string, 0, len(fieldSet))
	for field := range fieldSet {
		fields = append(fields, field)
	}

	return fields
}

// ReplaceVariables substitutes [placeholders] với actual values
func ReplaceVariables(text string, values map[string]string) string {
	re := regexp.MustCompile(`\[([^\]]+)\]`)

	return re.ReplaceAllStringFunc(text, func(match string) string {
		// Extract field name: "[barcode]" → "barcode"
		fieldName := match[1 : len(match)-1]

		// Replace với value, or keep placeholder if missing
		if value, ok := values[fieldName]; ok {
			return value
		}
		return match // Keep original if no value
	})
}
```

---

## Task 2.3: Implement Template Service

Create `internal/template/service.go`:

```go
package template

import (
	"vibe-imageborder/internal/models"
)

// Service handles template operations
type Service struct{}

// NewService creates a new TemplateService
func NewService() *Service {
	return &Service{}
}

// Load loads và validates template file
func (s *Service) Load(path string) (models.Template, error) {
	return ParseTemplate(path)
}

// GetDynamicFields extracts field names từ template
func (s *Service) GetDynamicFields(tmpl models.Template) []string {
	return ExtractDynamicFields(tmpl)
}

// ApplyValues replaces placeholders trong all template fields
func (s *Service) ApplyValues(tmpl models.Template, values map[string]string) models.Template {
	result := make(models.Template)

	for name, field := range tmpl {
		field.Text = ReplaceVariables(field.Text, values)
		result[name] = field
	}

	return result
}
```

---

## Task 2.4: Unit Tests

Create `internal/template/parser_test.go`:

```go
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
```

Run tests:

```bash
go test ./internal/template/...
```

---

## Task 2.5: Integration với Reference Templates

Copy reference templates từ UploadImage project:

```bash
# Copy test templates
cp "D:/Code-Tool/Software/web-tool/UploadImage/file/khung-002-05.txt" tests/fixtures/templates/
cp "D:/Code-Tool/Software/web-tool/UploadImage/file/khung-004-01.txt" tests/fixtures/templates/
```

Test với real templates:

```go
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
```

---

## Acceptance Criteria

- ✅ `ParseTemplate()` loads JSON correctly
- ✅ `validateTemplate()` catches invalid templates
- ✅ `ExtractDynamicFields()` finds all `[field]` placeholders
- ✅ `ReplaceVariables()` substitutes values correctly
- ✅ Unit tests pass (coverage >80%) - **Achieved 95.7%**
- ✅ Reference templates parse successfully - **khung-002-05.txt (7 fields), khung-004-01.txt (5 fields)**

**Code Review:** All critical, high, and medium priority items reviewed. No blocking issues. Ready for Phase 3.

---

## Deliverables

### Files Created

1. `internal/models/types.go` - Type definitions
2. `internal/template/parser.go` - Parsing logic
3. `internal/template/service.go` - Service wrapper
4. `internal/template/parser_test.go` - Unit tests
5. `tests/fixtures/templates/*.txt` - Test templates

### Validation Commands

```bash
# Run tests
go test ./internal/template/... -v -cover -race

# Actual output (2026-01-07)
=== RUN   TestParseTemplate
--- PASS: TestParseTemplate (0.00s)
=== RUN   TestExtractDynamicFields
--- PASS: TestExtractDynamicFields (0.00s)
=== RUN   TestReplaceVariables
--- PASS: TestReplaceVariables (0.00s)
=== RUN   TestRealTemplates
--- PASS: TestRealTemplates (0.00s)
=== RUN   TestParseTemplateErrors
--- PASS: TestParseTemplateErrors (0.01s)
=== RUN   TestReplaceVariablesEdgeCases
--- PASS: TestReplaceVariablesEdgeCases (0.00s)
=== RUN   TestService
--- PASS: TestService (0.00s)
PASS
coverage: 95.7% of statements
```

**Test Summary:**
- 7 test functions, 13 sub-tests
- All tests PASS
- Race detector: PASS
- Coverage: 95.7% (exceeds 80% target)

---

## Next Phase

[Phase 3: Image Service - Core Processing](phase-03-image-service-core.md)
