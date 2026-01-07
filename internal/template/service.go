package template

import (
	"vibe-imageborder/internal/models"
)

// Service handles template operations
type Service struct{}

// NewService creates a new TemplateService
func NewService() *Service {
	return &Service{}
}

// Load loads and validates template file
func (s *Service) Load(path string) (models.Template, error) {
	return ParseTemplate(path)
}

// GetDynamicFields extracts field names from template
func (s *Service) GetDynamicFields(tmpl models.Template) []string {
	return ExtractDynamicFields(tmpl)
}

// ApplyValues replaces placeholders in all template fields
func (s *Service) ApplyValues(tmpl models.Template, values map[string]string) models.Template {
	result := make(models.Template)

	for name, field := range tmpl {
		field.Text = ReplaceVariables(field.Text, values)
		result[name] = field
	}

	return result
}
