// Package image provides text rendering operations.
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

// namedColors maps color names to color values.
var namedColors = map[string]color.Color{
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

// TextRenderer handles text drawing.
type TextRenderer struct {
	fontManager *FontManager
}

// NewTextRenderer creates new text renderer.
func NewTextRenderer(fm *FontManager) *TextRenderer {
	return &TextRenderer{fontManager: fm}
}

// DrawOverlays draws all text overlays on image.
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

// drawSingleOverlay draws one text overlay.
func (tr *TextRenderer) drawSingleOverlay(dc *gg.Context, overlay models.TextOverlay) error {
	// Skip empty text
	if strings.TrimSpace(overlay.Text) == "" {
		return nil
	}

	// Parse position
	x, y, err := ParsePosition(overlay.Position)
	if err != nil {
		return fmt.Errorf("invalid position: %w", err)
	}

	// Parse color
	c := ParseColorName(overlay.Color)

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

// ParsePosition parses "x,y" string to coordinates.
func ParsePosition(pos string) (int, int, error) {
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

// ParseColorName converts color name or hex to color.Color.
func ParseColorName(name string) color.Color {
	name = strings.ToLower(strings.TrimSpace(name))

	if c, ok := namedColors[name]; ok {
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
