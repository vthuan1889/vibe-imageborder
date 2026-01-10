// Package template provides template operations.
package template

import (
	"vibe-imageborder/internal/models"
)

// Service handles template operations.
type Service struct{}

// NewService creates new template service.
func NewService() *Service {
	return &Service{}
}

// LoadTemplate loads template from file (always fresh read).
func (s *Service) LoadTemplate(path string) (*models.TemplateConfig, error) {
	return ParseTemplate(path)
}

// GetFields returns unique field names from template.
func (s *Service) GetFields(path string) ([]string, error) {
	config, err := s.LoadTemplate(path)
	if err != nil {
		return nil, err
	}
	return ExtractFields(config), nil
}

// GetOverlays returns text overlays with values applied.
func (s *Service) GetOverlays(path string, values map[string]string) (map[string]models.TextOverlay, error) {
	config, err := s.LoadTemplate(path)
	if err != nil {
		return nil, err
	}
	return ApplyValues(config, values), nil
}

// GetBackground returns background color from template.
func (s *Service) GetBackground(path string) (string, error) {
	config, err := s.LoadTemplate(path)
	if err != nil {
		return "", err
	}
	return config.Background, nil
}
