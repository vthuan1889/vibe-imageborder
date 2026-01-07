package models

import (
	"encoding/json"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

// Placeholder to ensure dependencies are imported
var _ = imaging.Open
var _ = gg.NewContext

// TemplateField represents a single text field in template
type TemplateField struct {
	Text     string `json:"text"`     // Text with [placeholders]
	Position string `json:"position"` // "x,y" format
	FontSize string `json:"fontsize"` // "45"
	Color    string `json:"color"`    // "white", "#FFFFFF"
}

// Template is a map of field name → field config
type Template map[string]TemplateField

// UnmarshalJSON custom parser to skip non-TemplateField entries (e.g., "background": "#f1eeea")
func (t *Template) UnmarshalJSON(data []byte) error {
	// Parse as raw map first
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = make(Template)
	for key, value := range raw {
		var field TemplateField
		// Try to unmarshal as TemplateField
		if err := json.Unmarshal(value, &field); err != nil {
			// If fails, skip (it's metadata like "background": "#f1eeea")
			continue
		}
		(*t)[key] = field
	}

	return nil
}

// CompositeRequest contains all data for batch processing
type CompositeRequest struct {
	ProductPaths []string          `json:"productPaths"` // Paths to product images
	FramePath    string            `json:"framePath"`    // Path to frame image
	TemplatePath string            `json:"templatePath"` // Path to template JSON
	FieldValues  map[string]string `json:"fieldValues"`  // User input: "barcode" → "ABC123"
	OutputDir    string            `json:"outputDir"`    // Output directory
}
