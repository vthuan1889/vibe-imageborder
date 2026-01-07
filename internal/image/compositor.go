package image

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"vibe-imageborder/internal/models"
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

// replaceVariables replaces [field] với values from map
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

// GetFontPath returns path to embedded font
func GetFontPath() string {
	return "assets/fonts/Roboto-Regular.ttf"
}

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

