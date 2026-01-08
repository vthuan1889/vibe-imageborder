// Package template provides template parsing and field extraction.
package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"vibe-imageborder/internal/models"
)

// fieldRegex matches placeholders like [field_name].
var fieldRegex = regexp.MustCompile(`\[([^\]]+)\]`)

// defaultFontSize is used when fontsize parsing fails.
const defaultFontSize = 24

// ParseTemplate reads and parses template JSON file.
func ParseTemplate(path string) (*models.TemplateConfig, error) {
	// Clean path to prevent directory traversal
	cleanPath := filepath.Clean(path)

	data, err := os.ReadFile(cleanPath)
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

// parseOverlay converts raw map to TextOverlay.
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
		size, err := strconv.Atoi(fs)
		if err != nil || size <= 0 {
			overlay.FontSize = defaultFontSize
		} else {
			overlay.FontSize = size
		}
	} else {
		overlay.FontSize = defaultFontSize
	}

	if color, ok := m["color"].(string); ok {
		overlay.Color = color
	}

	return overlay, nil
}

// ExtractFields returns unique field names from template.
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

// ApplyValues replaces placeholders with actual values.
// Skips overlays that have unfilled placeholders.
func ApplyValues(config *models.TemplateConfig, values map[string]string) map[string]models.TextOverlay {
	result := make(map[string]models.TextOverlay)

	for key, overlay := range config.Fields {
		newOverlay := overlay
		newOverlay.Text = replacePlaceholders(overlay.Text, values)
		
		// Skip if text still contains unfilled placeholders like [price]
		if !fieldRegex.MatchString(newOverlay.Text) {
			result[key] = newOverlay
		}
	}

	return result
}

// replacePlaceholders substitutes [field] with values.
func replacePlaceholders(text string, values map[string]string) string {
	result := text
	for field, value := range values {
		placeholder := "[" + field + "]"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
