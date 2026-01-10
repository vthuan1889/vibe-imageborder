// Package models defines shared data types for the image border application.
package models

// TextOverlay represents text to draw on image.
type TextOverlay struct {
	Text     string `json:"text"`
	Position string `json:"position"` // format: "x,y"
	FontSize int    `json:"fontsize"`
	Color    string `json:"color"`
}

// TemplateConfig represents parsed template configuration.
type TemplateConfig struct {
	Background string                 `json:"background,omitempty"`
	Fields     map[string]TextOverlay `json:"-"`
	FieldOrder []string               `json:"-"` // Preserves field order from JSON
	Raw        map[string]interface{} `json:"-"`
}

// ProcessRequest represents batch processing request from frontend.
type ProcessRequest struct {
	ProductImages []string          `json:"productImages"`
	FrameImage    string            `json:"frameImage"`
	TemplatePath  string            `json:"templatePath"`
	FieldValues   map[string]string `json:"fieldValues"`
	OutputDir     string            `json:"outputDir"`
	Format        string            `json:"format"` // png, jpg, webp
	Quality       int               `json:"quality"`
}

// ProcessProgress represents progress update during batch processing.
type ProcessProgress struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	File    string `json:"file"`
	Success bool   `json:"success"`
}
