# Phase 2: Template Service

## Context

- Plan: [plan.md](./plan.md)
- Previous: [Phase 1 - Project Setup](./phase-01-project-setup.md)

## Overview

| Field | Value |
|-------|-------|
| Priority | P1 - Critical Path |
| Status | Pending |
| Effort | 2h |

Implement template parsing service to read JSON template files and extract placeholder fields like `[barcode]`, `[price]`, `[size_dai]`.

## Requirements

### Functional
- Parse JSON template files (.txt extension containing JSON)
- Extract all `[field_name]` placeholders from text properties
- Return list of unique field names for dynamic form
- Validate template JSON structure
- Replace placeholders with actual values

### Non-functional
- Handle malformed JSON gracefully
- Support nested placeholders in single text property

## Architecture

```
Template Service
├── parser.go        # JSON parsing + validation
└── service.go       # Field extraction + value substitution

Flow:
.txt file → ParseTemplate() → TemplateConfig
                  ↓
         ExtractFields() → []string (unique field names)
                  ↓
         ApplyValues() → map[string]TextOverlay (with values)
```

## Template Format Reference

```json
{
  "background": "#f1eeea",
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
  },
  "size": {
    "text": "D[size_dai] x R[size_rong] x C[size_cao] CM",
    "position": "1100,10",
    "fontsize": "60",
    "color": "white"
  }
}
```

Expected extracted fields: `barcode`, `price`, `size_dai`, `size_rong`, `size_cao`

## Related Code Files

### Files to Create
| File | Purpose |
|------|---------|
| `internal/template/parser.go` | JSON parsing and validation |
| `internal/template/service.go` | Field extraction and value substitution |

### Files to Modify
| File | Change |
|------|--------|
| `internal/models/types.go` | Add any missing types |

## Implementation Steps

### Step 1: Create parser.go

```go
// internal/template/parser.go
package template

import (
    "encoding/json"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"

    "vibe-imageborder/internal/models"
)

// fieldRegex matches placeholders like [field_name]
var fieldRegex = regexp.MustCompile(`\[([^\]]+)\]`)

// ParseTemplate reads and parses template JSON file
func ParseTemplate(path string) (*models.TemplateConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read template: %w", err)
    }

    var raw map[string]interface{}
    if err := json.Unmarshal(data, &raw); err != nil {
        return nil, fmt.Errorf("invalid JSON: %w", err)
    }

    config := &models.TemplateConfig{
        Fields: make(map[string]models.TextOverlay),
        Raw:    raw,
    }

    // Extract background if exists
    if bg, ok := raw["background"].(string); ok {
        config.Background = bg
    }

    // Parse each field that has text overlay structure
    for key, val := range raw {
        if key == "background" {
            continue
        }

        overlay, err := parseOverlay(val)
        if err != nil {
            continue // Skip non-overlay fields
        }
        config.Fields[key] = overlay
    }

    return config, nil
}

// parseOverlay converts raw map to TextOverlay
func parseOverlay(val interface{}) (models.TextOverlay, error) {
    m, ok := val.(map[string]interface{})
    if !ok {
        return models.TextOverlay{}, fmt.Errorf("not a map")
    }

    overlay := models.TextOverlay{}

    if text, ok := m["text"].(string); ok {
        overlay.Text = text
    } else {
        return overlay, fmt.Errorf("missing text field")
    }

    if pos, ok := m["position"].(string); ok {
        overlay.Position = pos
    }

    if fs, ok := m["fontsize"].(string); ok {
        size, _ := strconv.Atoi(fs)
        overlay.FontSize = size
    }

    if color, ok := m["color"].(string); ok {
        overlay.Color = color
    }

    return overlay, nil
}

// ExtractFields returns unique field names from template
func ExtractFields(config *models.TemplateConfig) []string {
    fieldSet := make(map[string]bool)

    for _, overlay := range config.Fields {
        matches := fieldRegex.FindAllStringSubmatch(overlay.Text, -1)
        for _, match := range matches {
            if len(match) > 1 {
                fieldSet[match[1]] = true
            }
        }
    }

    fields := make([]string, 0, len(fieldSet))
    for field := range fieldSet {
        fields = append(fields, field)
    }
    return fields
}

// ApplyValues replaces placeholders with actual values
func ApplyValues(config *models.TemplateConfig, values map[string]string) map[string]models.TextOverlay {
    result := make(map[string]models.TextOverlay)

    for key, overlay := range config.Fields {
        newOverlay := overlay
        newOverlay.Text = replacePlaceholders(overlay.Text, values)
        result[key] = newOverlay
    }

    return result
}

// replacePlaceholders substitutes [field] with values
func replacePlaceholders(text string, values map[string]string) string {
    result := text
    for field, value := range values {
        placeholder := "[" + field + "]"
        result = strings.ReplaceAll(result, placeholder, value)
    }
    return result
}
```

### Step 2: Create service.go

```go
// internal/template/service.go
package template

import (
    "vibe-imageborder/internal/models"
)

// Service handles template operations
type Service struct {
    cache map[string]*models.TemplateConfig
}

// NewService creates new template service
func NewService() *Service {
    return &Service{
        cache: make(map[string]*models.TemplateConfig),
    }
}

// LoadTemplate loads and caches template
func (s *Service) LoadTemplate(path string) (*models.TemplateConfig, error) {
    // Check cache first
    if cached, ok := s.cache[path]; ok {
        return cached, nil
    }

    config, err := ParseTemplate(path)
    if err != nil {
        return nil, err
    }

    s.cache[path] = config
    return config, nil
}

// GetFields returns unique field names from template
func (s *Service) GetFields(path string) ([]string, error) {
    config, err := s.LoadTemplate(path)
    if err != nil {
        return nil, err
    }
    return ExtractFields(config), nil
}

// GetOverlays returns text overlays with values applied
func (s *Service) GetOverlays(path string, values map[string]string) (map[string]models.TextOverlay, error) {
    config, err := s.LoadTemplate(path)
    if err != nil {
        return nil, err
    }
    return ApplyValues(config, values), nil
}

// GetBackground returns background color from template
func (s *Service) GetBackground(path string) (string, error) {
    config, err := s.LoadTemplate(path)
    if err != nil {
        return "", err
    }
    return config.Background, nil
}

// ClearCache clears template cache
func (s *Service) ClearCache() {
    s.cache = make(map[string]*models.TemplateConfig)
}
```

### Step 3: Add Unit Tests

```go
// internal/template/parser_test.go
package template

import (
    "os"
    "path/filepath"
    "testing"
)

func TestExtractFields(t *testing.T) {
    // Create temp template file
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
    os.WriteFile(tmpFile, []byte(content), 0644)

    config, err := ParseTemplate(tmpFile)
    if err != nil {
        t.Fatalf("ParseTemplate failed: %v", err)
    }

    fields := ExtractFields(config)
    expected := map[string]bool{
        "barcode":   true,
        "size_dai":  true,
        "size_rong": true,
        "size_cao":  true,
    }

    if len(fields) != len(expected) {
        t.Errorf("Expected %d fields, got %d", len(expected), len(fields))
    }

    for _, f := range fields {
        if !expected[f] {
            t.Errorf("Unexpected field: %s", f)
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
        }
    }`

    tmpDir := t.TempDir()
    tmpFile := filepath.Join(tmpDir, "test.txt")
    os.WriteFile(tmpFile, []byte(content), 0644)

    config, _ := ParseTemplate(tmpFile)
    values := map[string]string{"barcode": "ABC123"}
    result := ApplyValues(config, values)

    if result["barcode"].Text != "ABC123" {
        t.Errorf("Expected ABC123, got %s", result["barcode"].Text)
    }
}
```

## Todo List

- [ ] Create `internal/template/parser.go`
- [ ] Create `internal/template/service.go`
- [ ] Create unit tests
- [ ] Test with real template files from reference app
- [ ] Verify field extraction works with complex templates

## Success Criteria

1. Parse template JSON without errors
2. Extract all unique field names correctly
3. Replace placeholders with values
4. Handle malformed JSON gracefully
5. Unit tests pass

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Template format variations | Medium | Flexible parsing, log unknown fields |
| Unicode in field names | Low | Go regex handles Unicode |
| Large template files | Low | Files are small (~1KB) |

## Security Considerations

- Validate file paths to prevent directory traversal
- Sanitize field values to prevent injection

## Next Steps

After completion, proceed to [Phase 3: Image Service Core](./phase-03-image-service-core.md)
