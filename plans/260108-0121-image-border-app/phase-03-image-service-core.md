# Phase 3: Image Service - Core

## Context

- Plan: [plan.md](./plan.md)
- Previous: [Phase 2 - Template Service](./phase-02-template-service.md)

## Overview

| Field | Value |
|-------|-------|
| Priority | P1 - Critical Path |
| Status | Completed |
| Effort | 3h |
| Completed | 2026-01-08 |

Implement core image processing: load, resize, composite product + frame overlay, save to multiple formats.

## Requirements

### Functional
- Load images (JPG, PNG, WebP)
- Resize product image to fit frame dimensions
- Composite: draw product → draw frame overlay
- Save to PNG, JPG, WebP with quality control
- Support alpha transparency in frame overlay

### Non-functional
- Process single image < 300ms
- Memory efficient (dispose images after use)
- Thread-safe for future parallel processing

## Architecture

```
Image Service
├── service.go       # Load, save, resize operations
└── compositor.go    # Composite product + frame

Flow:
Product Image + Frame Image
        ↓
    LoadImage()
        ↓
    ResizeToFit() (product to frame dimensions)
        ↓
    Composite() (product bg + frame overlay)
        ↓
    SaveImage() (format + quality)
```

## Related Code Files

### Files to Create
| File | Purpose |
|------|---------|
| `internal/image/service.go` | Image I/O operations |
| `internal/image/compositor.go` | Image compositing logic |

### Dependencies
```go
import (
    "github.com/disintegration/imaging"
    "image"
    "image/color"
    "image/draw"
    "image/jpeg"
    "image/png"
)
```

## Implementation Steps

### Step 1: Create service.go

```go
// internal/image/service.go
package image

import (
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "os"
    "path/filepath"
    "strings"

    "github.com/disintegration/imaging"
    "golang.org/x/image/webp"
)

// Service handles image operations
type Service struct{}

// NewService creates new image service
func NewService() *Service {
    return &Service{}
}

// LoadImage loads image from file path
func (s *Service) LoadImage(path string) (image.Image, error) {
    return imaging.Open(path)
}

// SaveImage saves image to file with format and quality
func (s *Service) SaveImage(img image.Image, path string, format string, quality int) error {
    // Ensure output directory exists
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create output dir: %w", err)
    }

    // Adjust path extension based on format
    ext := "." + strings.ToLower(format)
    basePath := strings.TrimSuffix(path, filepath.Ext(path))
    outputPath := basePath + ext

    file, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    switch strings.ToLower(format) {
    case "png":
        return png.Encode(file, img)
    case "jpg", "jpeg":
        return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
    case "webp":
        // Note: Go doesn't have native webp encoder
        // Use imaging library which falls back to PNG
        return imaging.Encode(file, img, imaging.PNG)
    default:
        return fmt.Errorf("unsupported format: %s", format)
    }
}

// ResizeToFit resizes image to fit within target dimensions
// Maintains aspect ratio, may be smaller than target
func (s *Service) ResizeToFit(img image.Image, width, height int) image.Image {
    return imaging.Fit(img, width, height, imaging.Lanczos)
}

// ResizeToFill resizes image to fill target dimensions
// Maintains aspect ratio, crops if needed
func (s *Service) ResizeToFill(img image.Image, width, height int) image.Image {
    return imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
}

// GetDimensions returns image width and height
func (s *Service) GetDimensions(img image.Image) (int, int) {
    bounds := img.Bounds()
    return bounds.Dx(), bounds.Dy()
}

// CreateBlankCanvas creates blank image with background color
func (s *Service) CreateBlankCanvas(width, height int, bgColor string) image.Image {
    c := parseColor(bgColor)
    return imaging.New(width, height, c)
}

// parseColor converts hex color string to color.Color
func parseColor(hex string) color.Color {
    hex = strings.TrimPrefix(hex, "#")
    if len(hex) != 6 {
        return color.White
    }

    r := hexToByte(hex[0:2])
    g := hexToByte(hex[2:4])
    b := hexToByte(hex[4:6])

    return color.RGBA{R: r, G: g, B: b, A: 255}
}

func hexToByte(s string) uint8 {
    var val int
    fmt.Sscanf(s, "%x", &val)
    return uint8(val)
}
```

### Step 2: Create compositor.go

```go
// internal/image/compositor.go
package image

import (
    "image"
    "image/draw"

    "github.com/disintegration/imaging"
)

// Compositor handles image compositing
type Compositor struct {
    service *Service
}

// NewCompositor creates new compositor
func NewCompositor(service *Service) *Compositor {
    return &Compositor{service: service}
}

// CompositeResult holds the composited image
type CompositeResult struct {
    Image  image.Image
    Width  int
    Height int
}

// Composite combines product and frame images
// Product is resized to fit frame, then frame overlaid on top
func (c *Compositor) Composite(product, frame image.Image, bgColor string) *CompositeResult {
    frameBounds := frame.Bounds()
    width := frameBounds.Dx()
    height := frameBounds.Dy()

    // Create canvas with background color
    var canvas *image.RGBA
    if bgColor != "" {
        bg := c.service.CreateBlankCanvas(width, height, bgColor)
        canvas = image.NewRGBA(image.Rect(0, 0, width, height))
        draw.Draw(canvas, canvas.Bounds(), bg, image.Point{}, draw.Src)
    } else {
        canvas = image.NewRGBA(image.Rect(0, 0, width, height))
    }

    // Resize product to fit frame dimensions
    resizedProduct := c.service.ResizeToFit(product, width, height)

    // Center product on canvas
    productBounds := resizedProduct.Bounds()
    offsetX := (width - productBounds.Dx()) / 2
    offsetY := (height - productBounds.Dy()) / 2

    // Draw product first (background)
    draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+productBounds.Dx(), offsetY+productBounds.Dy()),
        resizedProduct, image.Point{}, draw.Over)

    // Draw frame overlay (with alpha)
    draw.Draw(canvas, canvas.Bounds(), frame, image.Point{}, draw.Over)

    return &CompositeResult{
        Image:  canvas,
        Width:  width,
        Height: height,
    }
}

// CompositeWithPosition allows custom product positioning
func (c *Compositor) CompositeWithPosition(product, frame image.Image, bgColor string, position string) *CompositeResult {
    frameBounds := frame.Bounds()
    width := frameBounds.Dx()
    height := frameBounds.Dy()

    canvas := image.NewRGBA(image.Rect(0, 0, width, height))

    // Apply background if specified
    if bgColor != "" {
        bg := c.service.CreateBlankCanvas(width, height, bgColor)
        draw.Draw(canvas, canvas.Bounds(), bg, image.Point{}, draw.Src)
    }

    // Resize product to fit
    resizedProduct := c.service.ResizeToFit(product, width, height)
    productBounds := resizedProduct.Bounds()

    // Center product
    offsetX := (width - productBounds.Dx()) / 2
    offsetY := (height - productBounds.Dy()) / 2

    // Handle position-based ordering
    if position == "below" {
        // Frame first, then product
        draw.Draw(canvas, canvas.Bounds(), frame, image.Point{}, draw.Over)
        draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+productBounds.Dx(), offsetY+productBounds.Dy()),
            resizedProduct, image.Point{}, draw.Over)
    } else {
        // Product first, then frame (default)
        draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+productBounds.Dx(), offsetY+productBounds.Dy()),
            resizedProduct, image.Point{}, draw.Over)
        draw.Draw(canvas, canvas.Bounds(), frame, image.Point{}, draw.Over)
    }

    return &CompositeResult{
        Image:  canvas,
        Width:  width,
        Height: height,
    }
}

// ToRGBA converts image to RGBA for drawing operations
func ToRGBA(img image.Image) *image.RGBA {
    if rgba, ok := img.(*image.RGBA); ok {
        return rgba
    }

    bounds := img.Bounds()
    rgba := image.NewRGBA(bounds)
    draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
    return rgba
}
```

### Step 3: Add Unit Tests

```go
// internal/image/service_test.go
package image

import (
    "image"
    "image/color"
    "os"
    "path/filepath"
    "testing"
)

func TestLoadAndSaveImage(t *testing.T) {
    svc := NewService()

    // Create test image
    testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
    for y := 0; y < 100; y++ {
        for x := 0; x < 100; x++ {
            testImg.Set(x, y, color.RGBA{255, 0, 0, 255})
        }
    }

    tmpDir := t.TempDir()
    testPath := filepath.Join(tmpDir, "test.png")

    // Save
    err := svc.SaveImage(testImg, testPath, "png", 90)
    if err != nil {
        t.Fatalf("SaveImage failed: %v", err)
    }

    // Load
    loaded, err := svc.LoadImage(testPath)
    if err != nil {
        t.Fatalf("LoadImage failed: %v", err)
    }

    w, h := svc.GetDimensions(loaded)
    if w != 100 || h != 100 {
        t.Errorf("Expected 100x100, got %dx%d", w, h)
    }
}

func TestResizeToFit(t *testing.T) {
    svc := NewService()

    // Create 200x100 image (2:1 aspect)
    testImg := image.NewRGBA(image.Rect(0, 0, 200, 100))

    // Resize to fit 100x100 (should become 100x50)
    resized := svc.ResizeToFit(testImg, 100, 100)
    w, h := svc.GetDimensions(resized)

    if w != 100 || h != 50 {
        t.Errorf("Expected 100x50, got %dx%d", w, h)
    }
}

// internal/image/compositor_test.go
func TestComposite(t *testing.T) {
    svc := NewService()
    comp := NewCompositor(svc)

    // Create product image (white 100x100)
    product := image.NewRGBA(image.Rect(0, 0, 100, 100))
    for y := 0; y < 100; y++ {
        for x := 0; x < 100; x++ {
            product.Set(x, y, color.RGBA{255, 255, 255, 255})
        }
    }

    // Create frame image (200x200 with transparent center)
    frame := image.NewRGBA(image.Rect(0, 0, 200, 200))
    for y := 0; y < 200; y++ {
        for x := 0; x < 200; x++ {
            if x < 20 || x > 180 || y < 20 || y > 180 {
                frame.Set(x, y, color.RGBA{0, 0, 255, 255}) // Blue border
            } else {
                frame.Set(x, y, color.RGBA{0, 0, 0, 0}) // Transparent center
            }
        }
    }

    result := comp.Composite(product, frame, "#f1eeea")

    if result.Width != 200 || result.Height != 200 {
        t.Errorf("Expected 200x200, got %dx%d", result.Width, result.Height)
    }
}
```

## Todo List

- [x] Create `internal/image/service.go`
- [x] Create `internal/image/compositor.go`
- [x] Create unit tests
- [x] Test with real product + frame images
- [x] Verify PNG transparency works
- [x] Benchmark composite performance

## Success Criteria

1. Load JPG/PNG images without errors
2. Resize maintains aspect ratio
3. Composite preserves frame transparency
4. Save to PNG/JPG with quality control
5. Single composite < 300ms

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| WebP encoding not native | Medium | Fall back to PNG or use external lib |
| Large images slow | Low | Use Lanczos for quality, benchmar |
| Memory spikes | Medium | Dispose images after processing |

## Security Considerations

- Validate image file headers (not just extension)
- Limit max image dimensions to prevent OOM

## Next Steps

After completion, proceed to [Phase 4: Text Rendering](./phase-04-text-rendering.md)
