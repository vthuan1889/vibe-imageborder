# Phase 4: Image Service - Text Rendering

## Context

- Plan: [plan.md](./plan.md)
- Previous: [Phase 3 - Image Service Core](./phase-03-image-service-core.md)

## Overview

| Field | Value |
|-------|-------|
| Priority | P1 - Critical Path |
| Status | Completed |
| Effort | 2h |

Implement text rendering on images using bundled fonts. Support Vietnamese diacritics, custom colors, positions, and font sizes.

## Requirements

### Functional
- Load bundled TTF fonts from embedded assets
- Draw text at x,y position
- Support font size and color
- Render Vietnamese text with diacritics correctly
- Parse color from name ("white") or hex ("#ffffff")

### Non-functional
- Font loading is one-time, cached
- Text rendering < 50ms per overlay
- Support any Unicode characters

## Architecture

```
Text Rendering (extends image service)
├── fonts.go         # Font loading and caching
└── text.go          # Text drawing operations

Flow:
Load fonts once at startup
        ↓
For each TextOverlay:
    ParsePosition() → x, y
    ParseColor() → color.Color
    DrawText() → modify canvas
        ↓
Return modified image
```

## Related Code Files

### Files to Create
| File | Purpose |
|------|---------|
| `internal/image/fonts.go` | Font loading and management |
| `internal/image/text.go` | Text rendering operations |

### Files to Modify
| File | Change |
|------|--------|
| `internal/image/compositor.go` | Add DrawText method |
| `main.go` | Expose fonts embed |

### Dependencies
```go
import (
    "github.com/fogleman/gg"
    "golang.org/x/image/font/opentype"
)
```

## Implementation Steps

### Step 1: Create fonts.go

```go
// internal/image/fonts.go
package image

import (
    "embed"
    "fmt"
    "sync"

    "github.com/fogleman/gg"
    "golang.org/x/image/font"
    "golang.org/x/image/font/opentype"
)

// FontManager handles font loading and caching
type FontManager struct {
    fonts embed.FS
    cache map[string]*opentype.Font
    mu    sync.RWMutex
}

// NewFontManager creates font manager with embedded fonts
func NewFontManager(fontsFS embed.FS) *FontManager {
    return &FontManager{
        fonts: fontsFS,
        cache: make(map[string]*opentype.Font),
    }
}

// LoadFont loads font from embedded FS
func (fm *FontManager) LoadFont(name string) (*opentype.Font, error) {
    fm.mu.RLock()
    if cached, ok := fm.cache[name]; ok {
        fm.mu.RUnlock()
        return cached, nil
    }
    fm.mu.RUnlock()

    fm.mu.Lock()
    defer fm.mu.Unlock()

    // Double-check after acquiring write lock
    if cached, ok := fm.cache[name]; ok {
        return cached, nil
    }

    path := "assets/fonts/" + name + ".ttf"
    data, err := fm.fonts.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read font %s: %w", name, err)
    }

    f, err := opentype.Parse(data)
    if err != nil {
        return nil, fmt.Errorf("failed to parse font %s: %w", name, err)
    }

    fm.cache[name] = f
    return f, nil
}

// GetFace returns font.Face for given font and size
func (fm *FontManager) GetFace(name string, size float64) (font.Face, error) {
    f, err := fm.LoadFont(name)
    if err != nil {
        return nil, err
    }

    face, err := opentype.NewFace(f, &opentype.FaceOptions{
        Size:    size,
        DPI:     72,
        Hinting: font.HintingFull,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create face: %w", err)
    }

    return face, nil
}

// DefaultFontName returns default font to use
func DefaultFontName() string {
    return "BeVietnamPro-Regular"
}
```

### Step 2: Create text.go

```go
// internal/image/text.go
package image

import (
    "fmt"
    "image"
    "image/color"
    "strconv"
    "strings"

    "github.com/fogleman/gg"
    "vibe-imageborder/internal/models"
)

// TextRenderer handles text drawing
type TextRenderer struct {
    fontManager *FontManager
}

// NewTextRenderer creates new text renderer
func NewTextRenderer(fm *FontManager) *TextRenderer {
    return &TextRenderer{fontManager: fm}
}

// DrawOverlays draws all text overlays on image
func (tr *TextRenderer) DrawOverlays(img image.Image, overlays map[string]models.TextOverlay) (image.Image, error) {
    bounds := img.Bounds()
    dc := gg.NewContext(bounds.Dx(), bounds.Dy())
    dc.DrawImage(img, 0, 0)

    for _, overlay := range overlays {
        if err := tr.drawSingleOverlay(dc, overlay); err != nil {
            // Log error but continue with other overlays
            fmt.Printf("Warning: failed to draw overlay: %v\n", err)
            continue
        }
    }

    return dc.Image(), nil
}

// drawSingleOverlay draws one text overlay
func (tr *TextRenderer) drawSingleOverlay(dc *gg.Context, overlay models.TextOverlay) error {
    // Skip empty text
    if strings.TrimSpace(overlay.Text) == "" {
        return nil
    }

    // Parse position
    x, y, err := parsePosition(overlay.Position)
    if err != nil {
        return fmt.Errorf("invalid position: %w", err)
    }

    // Parse color
    c := parseColorName(overlay.Color)

    // Load font
    fontSize := float64(overlay.FontSize)
    if fontSize <= 0 {
        fontSize = 40 // default
    }

    face, err := tr.fontManager.GetFace(DefaultFontName(), fontSize)
    if err != nil {
        return fmt.Errorf("failed to load font: %w", err)
    }
    defer face.Close()

    dc.SetFontFace(face)
    dc.SetColor(c)
    dc.DrawString(overlay.Text, float64(x), float64(y)+fontSize) // gg uses baseline

    return nil
}

// parsePosition parses "x,y" string to coordinates
func parsePosition(pos string) (int, int, error) {
    parts := strings.Split(pos, ",")
    if len(parts) != 2 {
        return 0, 0, fmt.Errorf("invalid position format: %s", pos)
    }

    x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
    if err != nil {
        return 0, 0, fmt.Errorf("invalid x coordinate: %w", err)
    }

    y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
    if err != nil {
        return 0, 0, fmt.Errorf("invalid y coordinate: %w", err)
    }

    return x, y, nil
}

// parseColorName converts color name or hex to color.Color
func parseColorName(name string) color.Color {
    name = strings.ToLower(strings.TrimSpace(name))

    // Named colors
    colors := map[string]color.Color{
        "white":   color.White,
        "black":   color.Black,
        "red":     color.RGBA{255, 0, 0, 255},
        "green":   color.RGBA{0, 255, 0, 255},
        "blue":    color.RGBA{0, 0, 255, 255},
        "yellow":  color.RGBA{255, 255, 0, 255},
        "cyan":    color.RGBA{0, 255, 255, 255},
        "magenta": color.RGBA{255, 0, 255, 255},
        "gray":    color.RGBA{128, 128, 128, 255},
        "grey":    color.RGBA{128, 128, 128, 255},
    }

    if c, ok := colors[name]; ok {
        return c
    }

    // Try hex color
    if strings.HasPrefix(name, "#") {
        hex := strings.TrimPrefix(name, "#")
        if len(hex) == 6 {
            r := hexToByte(hex[0:2])
            g := hexToByte(hex[2:4])
            b := hexToByte(hex[4:6])
            return color.RGBA{R: r, G: g, B: b, A: 255}
        }
    }

    return color.White // default
}
```

### Step 3: Update compositor.go

Add text rendering to compositor:

```go
// Add to compositor.go

// CompositeWithText combines product, frame, and text overlays
func (c *Compositor) CompositeWithText(
    product, frame image.Image,
    bgColor string,
    overlays map[string]models.TextOverlay,
    textRenderer *TextRenderer,
) (*CompositeResult, error) {
    // First composite product + frame
    result := c.Composite(product, frame, bgColor)

    // Then draw text overlays
    if textRenderer != nil && len(overlays) > 0 {
        imgWithText, err := textRenderer.DrawOverlays(result.Image, overlays)
        if err != nil {
            return nil, fmt.Errorf("failed to draw text: %w", err)
        }
        result.Image = imgWithText
    }

    return result, nil
}
```

### Step 4: Update main.go for font embedding

```go
// In main.go, update embed directive:

//go:embed assets/fonts/*
var fonts embed.FS

// Pass to App:
func main() {
    app := NewApp(fonts)
    // ...
}
```

### Step 5: Add Unit Tests

```go
// internal/image/text_test.go
package image

import (
    "image"
    "image/color"
    "testing"

    "vibe-imageborder/internal/models"
)

func TestParsePosition(t *testing.T) {
    tests := []struct {
        input string
        x, y  int
        err   bool
    }{
        {"100,200", 100, 200, false},
        {"0,0", 0, 0, false},
        {"invalid", 0, 0, true},
        {"100", 0, 0, true},
    }

    for _, tt := range tests {
        x, y, err := parsePosition(tt.input)
        if tt.err && err == nil {
            t.Errorf("Expected error for %s", tt.input)
        }
        if !tt.err && (x != tt.x || y != tt.y) {
            t.Errorf("For %s: expected %d,%d got %d,%d", tt.input, tt.x, tt.y, x, y)
        }
    }
}

func TestParseColorName(t *testing.T) {
    tests := []struct {
        input    string
        expected color.Color
    }{
        {"white", color.White},
        {"WHITE", color.White},
        {"black", color.Black},
        {"#ff0000", color.RGBA{255, 0, 0, 255}},
        {"#00FF00", color.RGBA{0, 255, 0, 255}},
    }

    for _, tt := range tests {
        result := parseColorName(tt.input)
        r1, g1, b1, a1 := result.RGBA()
        r2, g2, b2, a2 := tt.expected.RGBA()
        if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
            t.Errorf("For %s: color mismatch", tt.input)
        }
    }
}
```

## Todo List

- [x] Create `internal/image/fonts.go`
- [x] Create `internal/image/text.go`
- [x] Update `compositor.go` with text method
- [x] Update `main.go` for font embedding
- [x] Add unit tests
- [x] Test Vietnamese text rendering
- [x] Verify font loading works from embedded FS

## Success Criteria

1. Load bundled fonts without external dependencies
2. Draw text at correct positions
3. Vietnamese diacritics render correctly
4. Color names and hex codes work
5. Text rendering < 50ms per overlay

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Font missing glyphs | High | Use Be Vietnam Pro (full Vietnamese) |
| Wrong text position | Medium | Test with real templates |
| Font embedding size | Low | ~200KB per font, acceptable |

## Security Considerations

- Fonts embedded, no external loading
- Sanitize text input (no control characters)

## Next Steps

After completion, proceed to [Phase 5: Wails Backend](./phase-05-wails-backend.md)
