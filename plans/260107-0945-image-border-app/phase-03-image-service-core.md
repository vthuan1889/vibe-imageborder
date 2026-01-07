# Phase 3: Image Service - Core Processing

**Goal:** Load images, resize với contain mode, composite (overlay) without text

**Duration:** ~3-4 hours

**Dependencies:** Phase 1-2 complete

---

## Overview

Implement core image operations:
1. Load JPEG/PNG images
2. Resize product images to fit frame (contain mode)
3. Composite product overlay lên frame (centered)
4. Save output images

**Note:** Text rendering sẽ được thêm trong Phase 4.

---

## Task 3.1: Implement Image Loading

Create `internal/image/service.go`:

```go
package image

import (
	"fmt"
	"image"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// Service handles image operations
type Service struct {
	// Future: Add configuration if needed
}

// NewService creates a new ImageService
func NewService() *Service {
	return &Service{}
}

// LoadImage loads an image from path (JPEG/PNG)
func (s *Service) LoadImage(path string) (image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load image %s: %w", filepath.Base(path), err)
	}

	// Validate image dimensions (prevent OOM)
	bounds := img.Bounds()
	maxDim := 10000
	if bounds.Dx() > maxDim || bounds.Dy() > maxDim {
		return nil, fmt.Errorf("image too large: %dx%d (max: %dx%d)",
			bounds.Dx(), bounds.Dy(), maxDim, maxDim)
	}

	return img, nil
}

// SaveImage saves an image to path as PNG
func (s *Service) SaveImage(img image.Image, path string) error {
	err := imaging.Save(img, path)
	if err != nil {
		return fmt.Errorf("failed to save image %s: %w", filepath.Base(path), err)
	}
	return nil
}
```

---

## Task 3.2: Implement Resize Logic

Add to `internal/image/service.go`:

```go
// ResizeToFit resizes image to fit within bounds (contain mode)
// Preserves aspect ratio, adds letterbox if needed
func (s *Service) ResizeToFit(img image.Image, maxWidth, maxHeight int) image.Image {
	// Use imaging.Fit for contain mode
	// Lanczos resampling for high quality
	return imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
}
```

**Rationale:**
- `imaging.Fit()` implements contain mode perfectly
- Lanczos filter provides best quality for downscaling
- Aspect ratio preserved automatically

---

## Task 3.3: Implement Image Compositing

Create `internal/image/compositor.go`:

```go
package image

import (
	"fmt"
	"image"
)

// CompositeResult contains output và metadata
type CompositeResult struct {
	Image        image.Image
	ProductPath  string
	Success      bool
	ErrorMessage string
}

// CompositeImages overlays product image lên frame (centered)
func (s *Service) CompositeImages(productPath, framePath string) (*CompositeResult, error) {
	result := &CompositeResult{
		ProductPath: productPath,
		Success:     false,
	}

	// Load product image
	productImg, err := s.LoadImage(productPath)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result, err
	}

	// Load frame image
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

	// Create composite (will use gg.Context later for text)
	// For now, use basic imaging.Overlay
	composite := imaging.Clone(frameImg)
	composite = imaging.Overlay(composite, productFit, image.Pt(centerX, centerY), 1.0)

	result.Image = composite
	result.Success = true

	return result, nil
}
```

**Note:** Phase 4 sẽ replace `imaging.Overlay` bằng `gg.Context` để support text rendering.

---

## Task 3.4: CLI Test Program

Create `cmd/test-composite/main.go`:

```go
package main

import (
	"fmt"
	"log"
	"path/filepath"
	"vibe-imageborder/internal/image"
)

func main() {
	fmt.Println("Image Compositing Test")
	fmt.Println("======================")

	svc := image.NewService()

	// Test paths
	productPath := "tests/fixtures/products/product-01.jpg"
	framePath := "tests/fixtures/frames/frame-01.png"
	outputPath := "tests/output/composite-test.png"

	// Composite
	result, err := svc.CompositeImages(productPath, framePath)
	if err != nil {
		log.Fatalf("Composite failed: %v", err)
	}

	if !result.Success {
		log.Fatalf("Composite unsuccessful: %s", result.ErrorMessage)
	}

	// Save output
	err = svc.SaveImage(result.Image, outputPath)
	if err != nil {
		log.Fatalf("Save failed: %v", err)
	}

	fmt.Printf("✓ Success! Output saved to: %s\n", outputPath)
}
```

Run test:

```bash
# Create output directory
mkdir -p tests/output

# Run test
go run cmd/test-composite/main.go
```

**Expected:**
- Output image created tại `tests/output/composite-test.png`
- Product image centered trong frame
- No errors

---

## Task 3.5: Unit Tests

Create `internal/image/service_test.go`:

```go
package image

import (
	"image"
	"image/color"
	"testing"

	"github.com/disintegration/imaging"
)

func TestLoadImage(t *testing.T) {
	svc := NewService()

	// Create test image
	tmpImg := imaging.New(800, 600, color.RGBA{255, 0, 0, 255})
	tmpPath := "../../tests/output/test_load.png"
	imaging.Save(tmpImg, tmpPath)

	// Test load
	img, err := svc.LoadImage(tmpPath)
	if err != nil {
		t.Fatalf("LoadImage failed: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 800 || bounds.Dy() != 600 {
		t.Errorf("Expected 800x600, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestResizeToFit(t *testing.T) {
	svc := NewService()

	// Create test image: 1000x500 (2:1 aspect ratio)
	img := imaging.New(1000, 500, color.RGBA{0, 255, 0, 255})

	// Fit into 800x800 (should be 800x400 to preserve 2:1)
	resized := svc.ResizeToFit(img, 800, 800)

	bounds := resized.Bounds()
	expectedW, expectedH := 800, 400

	if bounds.Dx() != expectedW || bounds.Dy() != expectedH {
		t.Errorf("Expected %dx%d, got %dx%d",
			expectedW, expectedH, bounds.Dx(), bounds.Dy())
	}
}

func TestCompositeImages(t *testing.T) {
	svc := NewService()

	// Create test images
	frame := imaging.New(1000, 1000, color.RGBA{200, 200, 200, 255})
	product := imaging.New(600, 400, color.RGBA{255, 0, 0, 255})

	framePath := "../../tests/output/test_frame.png"
	productPath := "../../tests/output/test_product.png"

	imaging.Save(frame, framePath)
	imaging.Save(product, productPath)

	// Test composite
	result, err := svc.CompositeImages(productPath, framePath)
	if err != nil {
		t.Fatalf("CompositeImages failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Composite unsuccessful: %s", result.ErrorMessage)
	}

	// Verify output dimensions match frame
	bounds := result.Image.Bounds()
	if bounds.Dx() != 1000 || bounds.Dy() != 1000 {
		t.Errorf("Expected 1000x1000, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}
```

Run tests:

```bash
go test ./internal/image/... -v -cover
```

---

## Task 3.6: Add Test Fixtures

Prepare test images:

```bash
# Create simple test frames và products
# Option 1: Download sample images
# Option 2: Use ImageMagick to generate

# Generate test frame (1000x1000 gray border)
convert -size 1000x1000 xc:white \
  -fill none -stroke gray -strokewidth 50 \
  -draw "rectangle 25,25 975,975" \
  tests/fixtures/frames/frame-01.png

# Generate test product (600x400 red rectangle)
convert -size 600x400 xc:red \
  tests/fixtures/products/product-01.jpg
```

**Alternative:** Copy sample images từ reference project.

---

## Acceptance Criteria

- ✓ `LoadImage()` loads JPEG/PNG correctly
- ✓ `LoadImage()` validates image size (max 10000x10000)
- ✓ `ResizeToFit()` preserves aspect ratio
- ✓ `CompositeImages()` centers product trong frame
- ✓ `SaveImage()` outputs PNG correctly
- ✓ Unit tests pass với >80% coverage
- ✓ CLI test program works end-to-end

---

## Deliverables

### Files Created

1. `internal/image/service.go` - Image operations
2. `internal/image/compositor.go` - Compositing logic
3. `internal/image/service_test.go` - Unit tests
4. `cmd/test-composite/main.go` - CLI test program
5. `tests/fixtures/frames/*.png` - Test frames
6. `tests/fixtures/products/*.jpg` - Test products

### Validation

```bash
# 1. Run unit tests
go test ./internal/image/... -v -cover

# 2. Run CLI test
go run cmd/test-composite/main.go

# 3. Verify output
open tests/output/composite-test.png  # macOS
# or
start tests/output/composite-test.png  # Windows
```

**Visual Check:**
- Product image centered
- Aspect ratio preserved
- No distortion or quality loss

---

## Known Limitations

- Text rendering not implemented (Phase 4)
- Single-threaded processing (Phase 8 enhancement)
- Fixed PNG output format

These will be addressed in subsequent phases.

---

## Next Phase

[Phase 4: Image Service - Text Rendering](phase-04-text-rendering.md)
