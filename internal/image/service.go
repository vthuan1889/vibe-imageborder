package image

import (
	"fmt"
	"image"
	"path/filepath"

	"github.com/disintegration/imaging"
	"vibe-imageborder/internal/models"
)

// Service handles image operations
type Service struct {
	// Future: Add configuration if needed
}

// NewService creates a new ImageService
func NewService() *Service {
	return &Service{}
}

// LoadImage loads an image from path (JPEG/PNG)
func (s *Service) LoadImage(path string) (image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load image %s: %w", filepath.Base(path), err)
	}

	// Validate image dimensions (prevent OOM)
	bounds := img.Bounds()
	maxDim := 10000
	if bounds.Dx() > maxDim || bounds.Dy() > maxDim {
		return nil, fmt.Errorf("image too large: %dx%d (max: %dx%d)",
			bounds.Dx(), bounds.Dy(), maxDim, maxDim)
	}

	return img, nil
}

// SaveImage saves an image to path as PNG
func (s *Service) SaveImage(img image.Image, path string) error {
	err := imaging.Save(img, path)
	if err != nil {
		return fmt.Errorf("failed to save image %s: %w", filepath.Base(path), err)
	}
	return nil
}

// ResizeToFit resizes image to fit within bounds (contain mode)
// Preserves aspect ratio, adds letterbox if needed
func (s *Service) ResizeToFit(img image.Image, maxWidth, maxHeight int) image.Image {
	// Use imaging.Fit for contain mode
	// Lanczos resampling for high quality
	return imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
}

// ProcessSingle processes one product image với template
func (s *Service) ProcessSingle(
	productPath, framePath, outputPath string,
	template models.Template,
	fieldValues map[string]string,
) error {
	result, err := s.CompositeImagesWithText(productPath, framePath, template, fieldValues)
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("composite failed: %s", result.ErrorMessage)
	}

	return s.SaveImage(result.Image, outputPath)
}

