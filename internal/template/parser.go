package template

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"vibe-imageborder/internal/models"
)

// ParseTemplate loads and parses template JSON file
func ParseTemplate(path string) (models.Template, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Parse JSON
	var tmpl models.Template
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse template JSON: %w", err)
	}

	// Validate template
	if err := validateTemplate(tmpl); err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}

	return tmpl, nil
}

// validateTemplate checks template structure
func validateTemplate(tmpl models.Template) error {
	if len(tmpl) == 0 {
		return fmt.Errorf("template is empty")
	}

	for name, field := range tmpl {
		if field.Text == "" {
			return fmt.Errorf("field %s: text is empty", name)
		}
		if field.Position == "" {
			return fmt.Errorf("field %s: position is empty", name)
		}
		if field.FontSize == "" {
			return fmt.Errorf("field %s: fontsize is empty", name)
		}
		if field.Color == "" {
			return fmt.Errorf("field %s: color is empty", name)
		}
	}

	return nil
}

// ExtractDynamicFields finds all unique [field] placeholders in template
func ExtractDynamicFields(tmpl models.Template) []string {
	fieldSet := make(map[string]bool)
	re := regexp.MustCompile(`\[([^\]]+)\]`)

	for _, field := range tmpl {
		matches := re.FindAllStringSubmatch(field.Text, -1)
		for _, match := range matches {
			if len(match) > 1 {
				fieldSet[match[1]] = true
			}
		}
	}

	// Convert set to slice
	fields := make([]string, 0, len(fieldSet))
	for field := range fieldSet {
		fields = append(fields, field)
	}

	return fields
}

// ReplaceVariables substitutes [placeholders] with actual values
func ReplaceVariables(text string, values map[string]string) string {
	re := regexp.MustCompile(`\[([^\]]+)\]`)

	return re.ReplaceAllStringFunc(text, func(match string) string {
		// Extract field name: "[barcode]" → "barcode"
		fieldName := match[1 : len(match)-1]

		// Replace with value, or keep placeholder if missing
		if value, ok := values[fieldName]; ok {
			return value
		}
		return match // Keep original if no value
	})
}
