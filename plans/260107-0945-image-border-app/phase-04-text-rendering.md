# Phase 4: Image Service - Text Rendering

**Goal:** Add text overlay với TTF fonts, position-based placement

**Duration:** ~3-4 hours

**Dependencies:** Phase 1-3 complete

---

## Overview

Upgrade image compositing để support text rendering:
1. Embed TTF font trong app
2. Parse color strings → RGBA
3. Parse position strings → (x,y)
4. Render text với `fogleman/gg`
5. Replace `imaging.Overlay` với `gg.Context` workflow

---

## Task 4.1: Embed Font Assets

Download Roboto font:

```bash
# Download Roboto-Regular.ttf
curl -L -o assets/fonts/Roboto-Regular.ttf \
  "https://github.com/google/fonts/raw/main/apache/roboto/static/Roboto-Regular.ttf"
```

Update `app.go` để embed font:

```go
package main

import (
	"embed"
	// ... other imports
)

//go:embed assets/fonts/*
var fontsFS embed.FS

// GetFontPath returns path to embedded font
func GetFontPath() string {
	// Font will be extracted to temp at runtime
	// For now, use relative path (Wails handles this)
	return "assets/fonts/Roboto-Regular.ttf"
}
```

---

## Task 4.2: Implement Color Parsing

Add to `internal/image/compositor.go`:

```go
import (
	"image/color"
	"strconv"
	"strings"
)

// parseColor converts color string to color.Color
// Supports: "white", "black", "red", "#RRGGBB", "#RRGGBBAA"
func parseColor(colorStr string) color.Color {
	colorStr = strings.TrimSpace(strings.ToLower(colorStr))

	// Named colors
	switch colorStr {
	case "white":
		return color.RGBA{255, 255, 255, 255}
	case "black":
		return color.RGBA{0, 0, 0, 255}
	case "red":
		return color.RGBA{255, 0, 0, 255}
	case "green":
		return color.RGBA{0, 255, 0, 255}
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	case "yellow":
		return color.RGBA{255, 255, 0, 255}
	}

	// Hex colors: #RGB, #RRGGBB, #RRGGBBAA
	if strings.HasPrefix(colorStr, "#") {
		hex := colorStr[1:]

		// Expand #RGB to #RRGGBB
		if len(hex) == 3 {
			hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
		}

		if len(hex) == 6 {
			hex += "FF" // Add full alpha
		}

		if len(hex) == 8 {
			r, _ := strconv.ParseUint(hex[0:2], 16, 8)
			g, _ := strconv.ParseUint(hex[2:4], 16, 8)
			b, _ := strconv.ParseUint(hex[4:6], 16, 8)
			a, _ := strconv.ParseUint(hex[6:8], 16, 8)
			return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		}
	}

	// Default: white
	return color.RGBA{255, 255, 255, 255}
}
```

---

## Task 4.3: Implement Position Parsing

Add to `internal/image/compositor.go`:

```go
import (
	"fmt"
	"strconv"
)

// parsePosition converts "x,y" string to (x, y) coordinates
func parsePosition(posStr string) (int, int, error) {
	parts := strings.Split(posStr, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid position format: %s (expected \"x,y\")", posStr)
	}

	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid x coordinate: %s", parts[0])
	}

	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid y coordinate: %s", parts[1])
	}

	return x, y, nil
}

// parseFontSize converts fontsize string to float64
func parseFontSize(sizeStr string) (float64, error) {
	size, err := strconv.ParseFloat(strings.TrimSpace(sizeStr), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fontsize: %s", sizeStr)
	}
	return size, nil
}
```

---

## Task 4.4: Implement Text Rendering với gg

Refactor `CompositeImages()` trong `internal/image/compositor.go`:

```go
import (
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"vibe-imageborder/internal/models"
)

// CompositeImagesWithText overlays product + renders text fields
func (s *Service) CompositeImagesWithText(
	productPath, framePath string,
	template models.Template,
	fieldValues map[string]string,
) (*CompositeResult, error) {
	result := &CompositeResult{
		ProductPath: productPath,
		Success:     false,
	}

	// Load images
	productImg, err := s.LoadImage(productPath)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result, err
	}

	frameImg, err := s.LoadImage(framePath)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result, err
	}

	// Resize product to fit frame
	frameBounds := frameImg.Bounds()
	productFit := s.ResizeToFit(productImg, frameBounds.Dx(), frameBounds.Dy())

	// Calculate center position
	productBounds := productFit.Bounds()
	centerX := (frameBounds.Dx() - productBounds.Dx()) / 2
	centerY := (frameBounds.Dy() - productBounds.Dy()) / 2

	// Create gg.Context from frame
	dc := gg.NewContextForImage(frameImg)

	// Draw product image centered
	dc.DrawImage(productFit, centerX, centerY)

	// Render text fields
	for fieldName, field := range template {
		// Replace variables trong text
		text := replaceVariables(field.Text, fieldValues)

		// Parse field properties
		x, y, err := parsePosition(field.Position)
		if err != nil {
			log.Printf("Warning: %v for field %s", err, fieldName)
			continue
		}

		fontSize, err := parseFontSize(field.FontSize)
		if err != nil {
			log.Printf("Warning: %v for field %s", err, fieldName)
			continue
		}

		textColor := parseColor(field.Color)

		// Load font
		fontPath := GetFontPath()
		if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
			log.Printf("Warning: Failed to load font for field %s: %v", fieldName, err)
			continue
		}

		// Set color và draw text
		dc.SetColor(textColor)
		dc.DrawString(text, float64(x), float64(y))
	}

	result.Image = dc.Image()
	result.Success = true

	return result, nil
}

// replaceVariables helper (moved from template package for convenience)
func replaceVariables(text string, values map[string]string) string {
	re := regexp.MustCompile(`\[([^\]]+)\]`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		fieldName := match[1 : len(match)-1]
		if value, ok := values[fieldName]; ok {
			return value
		}
		return match
	})
}
```

---

## Task 4.5: Update Service Interface

Add to `internal/image/service.go`:

```go
// ProcessSingle processes one product image với template
func (s *Service) ProcessSingle(
	productPath, framePath, outputPath string,
	template models.Template,
	fieldValues map[string]string,
) error {
	result, err := s.CompositeImagesWithText(productPath, framePath, template, fieldValues)
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("composite failed: %s", result.ErrorMessage)
	}

	return s.SaveImage(result.Image, outputPath)
}
```

---

## Task 4.6: Update CLI Test Program

Update `cmd/test-composite/main.go`:

```go
package main

import (
	"fmt"
	"log"
	"vibe-imageborder/internal/image"
	"vibe-imageborder/internal/template"
)

func main() {
	fmt.Println("Image Compositing với Text Test")
	fmt.Println("=================================")

	imgSvc := image.NewService()
	tmplSvc := template.NewService()

	// Load template
	templatePath := "tests/fixtures/templates/khung-002-05.txt"
	tmpl, err := tmplSvc.Load(templatePath)
	if err != nil {
		log.Fatalf("Failed to load template: %v", err)
	}

	// Dynamic fields
	fields := tmplSvc.GetDynamicFields(tmpl)
	fmt.Printf("Template fields: %v\n", fields)

	// Field values
	fieldValues := map[string]string{
		"barcode":   "ABC123456",
		"size_dai":  "30",
		"size_rong": "20",
		"size_cao":  "15",
	}

	// Process
	productPath := "tests/fixtures/products/product-01.jpg"
	framePath := "tests/fixtures/frames/frame-01.png"
	outputPath := "tests/output/composite-with-text.png"

	err = imgSvc.ProcessSingle(productPath, framePath, outputPath, tmpl, fieldValues)
	if err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	fmt.Printf("✓ Success! Output saved to: %s\n", outputPath)
}
```

Run test:

```bash
go run cmd/test-composite/main.go
```

**Expected:**
- Output image với product centered
- Text fields rendered tại correct positions
- Colors và font sizes correct

---

## Task 4.7: Unit Tests

Add to `internal/image/compositor_test.go`:

```go
package image

import (
	"image/color"
	"testing"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		input    string
		expected color.Color
	}{
		{"white", color.RGBA{255, 255, 255, 255}},
		{"black", color.RGBA{0, 0, 0, 255}},
		{"#FFFFFF", color.RGBA{255, 255, 255, 255}},
		{"#FF0000", color.RGBA{255, 0, 0, 255}},
		{"#FF0000AA", color.RGBA{255, 0, 0, 170}},
		{"invalid", color.RGBA{255, 255, 255, 255}}, // Default white
	}

	for _, tt := range tests {
		result := parseColor(tt.input)
		if result != tt.expected {
			t.Errorf("parseColor(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestParsePosition(t *testing.T) {
	tests := []struct {
		input     string
		expectX   int
		expectY   int
		expectErr bool
	}{
		{"98,1720", 98, 1720, false},
		{"0,0", 0, 0, false},
		{"100, 200", 100, 200, false}, // With spaces
		{"invalid", 0, 0, true},
		{"100", 0, 0, true}, // Missing y
	}

	for _, tt := range tests {
		x, y, err := parsePosition(tt.input)
		if tt.expectErr {
			if err == nil {
				t.Errorf("parsePosition(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("parsePosition(%q) unexpected error: %v", tt.input, err)
			}
			if x != tt.expectX || y != tt.expectY {
				t.Errorf("parsePosition(%q) = (%d, %d), want (%d, %d)",
					tt.input, x, y, tt.expectX, tt.expectY)
			}
		}
	}
}

func TestParseFontSize(t *testing.T) {
	tests := []struct {
		input     string
		expected  float64
		expectErr bool
	}{
		{"45", 45.0, false},
		{"12.5", 12.5, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		size, err := parseFontSize(tt.input)
		if tt.expectErr {
			if err == nil {
				t.Errorf("parseFontSize(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("parseFontSize(%q) unexpected error: %v", tt.input, err)
			}
			if size != tt.expected {
				t.Errorf("parseFontSize(%q) = %f, want %f", tt.input, size, tt.expected)
			}
		}
	}
}
```

Run tests:

```bash
go test ./internal/image/... -v -cover
```

---

## Acceptance Criteria

- ✓ Roboto-Regular.ttf embedded trong app
- ✓ `parseColor()` handles named colors và hex
- ✓ `parsePosition()` parses "x,y" correctly
- ✓ `parseFontSize()` converts string to float
- ✓ `CompositeImagesWithText()` renders text correctly
- ✓ CLI test outputs image với visible text
- ✓ Unit tests pass với >80% coverage

---

## Deliverables

### Files Created/Modified

1. `assets/fonts/Roboto-Regular.ttf` - Embedded font
2. `app.go` - Font embedding directive
3. `internal/image/compositor.go` - Text rendering logic
4. `internal/image/service.go` - ProcessSingle method
5. `internal/image/compositor_test.go` - Unit tests
6. `cmd/test-composite/main.go` - Updated test

### Validation

```bash
# 1. Run unit tests
go test ./internal/image/... -v -cover

# 2. Run CLI test với real template
go run cmd/test-composite/main.go

# 3. Visual check
open tests/output/composite-with-text.png
```

**Visual Verification:**
- Product centered trong frame ✓
- Text rendered tại positions from template ✓
- Font size matches template ✓
- Color matches template (white text visible) ✓

---

## Troubleshooting

### Issue: Font not loading

**Solution:**
```go
// Ensure font path is correct
fontPath := "assets/fonts/Roboto-Regular.ttf"

// Check file exists
if _, err := os.Stat(fontPath); os.IsNotExist(err) {
    log.Fatal("Font file not found: ", fontPath)
}
```

### Issue: Text not visible

**Possible causes:**
1. Color same as background (e.g., white on white)
2. Position outside image bounds
3. Font size too small

**Debug:**
```go
log.Printf("Rendering text %q at (%d,%d) size=%f color=%v",
    text, x, y, fontSize, textColor)
```

---

## Next Phase

[Phase 5: Wails Backend Integration](phase-05-wails-backend.md)
