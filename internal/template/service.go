// Package template provides template operations.
package template

import (
	"sync"

	"vibe-imageborder/internal/models"
)

// Service handles template operations with caching.
type Service struct {
	cache map[string]*models.TemplateConfig
	mu    sync.RWMutex
}

// NewService creates new template service.
func NewService() *Service {
	return &Service{
		cache: make(map[string]*models.TemplateConfig),
	}
}

// LoadTemplate loads and caches template.
func (s *Service) LoadTemplate(path string) (*models.TemplateConfig, error) {
	s.mu.RLock()
	if cached, ok := s.cache[path]; ok {
		s.mu.RUnlock()
		return cached, nil
	}
	s.mu.RUnlock()

	config, err := ParseTemplate(path)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.cache[path] = config
	s.mu.Unlock()
	return config, nil
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

// ClearCache clears template cache.
func (s *Service) ClearCache() {
	s.mu.Lock()
	s.cache = make(map[string]*models.TemplateConfig)
	s.mu.Unlock()
}
