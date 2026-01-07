// Package image provides image loading, saving, and transformation operations.
package image

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	// Import for webp decoding support
	_ "golang.org/x/image/webp"
)

// Maximum allowed image dimensions to prevent OOM.
const (
	MaxImageWidth  = 8192
	MaxImageHeight = 8192
)

// ErrImageTooLarge is returned when image exceeds maximum dimensions.
var ErrImageTooLarge = errors.New("image exceeds maximum allowed dimensions")

// Service handles image operations.
type Service struct{}

// NewService creates new image service.
func NewService() *Service {
	return &Service{}
}

// LoadImage loads image from file path with size validation.
func (s *Service) LoadImage(path string) (image.Image, error) {
	cleanPath := filepath.Clean(path)

	img, err := imaging.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}

	// Validate dimensions to prevent OOM
	bounds := img.Bounds()
	if bounds.Dx() > MaxImageWidth || bounds.Dy() > MaxImageHeight {
		return nil, fmt.Errorf("%w: %dx%d exceeds %dx%d",
			ErrImageTooLarge, bounds.Dx(), bounds.Dy(), MaxImageWidth, MaxImageHeight)
	}

	return img, nil
}

// SaveImage saves image to file with format and quality.
func (s *Service) SaveImage(img image.Image, path string, format string, quality int) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	ext := "." + strings.ToLower(format)
	basePath := strings.TrimSuffix(path, filepath.Ext(path))
	outputPath := basePath + ext

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	switch strings.ToLower(format) {
	case "png":
		return png.Encode(file, img)
	case "jpg", "jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	case "webp":
		// Go doesn't have native webp encoder, fall back to PNG
		return imaging.Encode(file, img, imaging.PNG)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// ResizeToFit resizes image to fit within target dimensions.
// Maintains aspect ratio, may be smaller than target.
func (s *Service) ResizeToFit(img image.Image, width, height int) image.Image {
	return imaging.Fit(img, width, height, imaging.Lanczos)
}

// ResizeToFill resizes image to fill target dimensions.
// Maintains aspect ratio, crops if needed.
func (s *Service) ResizeToFill(img image.Image, width, height int) image.Image {
	return imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
}

// GetDimensions returns image width and height.
func (s *Service) GetDimensions(img image.Image) (int, int) {
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

// CreateBlankCanvas creates blank image with background color.
func (s *Service) CreateBlankCanvas(width, height int, bgColor string) image.Image {
	c := parseColor(bgColor)
	return imaging.New(width, height, c)
}

// parseColor converts hex color string to color.Color.
func parseColor(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.White
	}

	r := hexToByte(hex[0:2])
	g := hexToByte(hex[2:4])
	b := hexToByte(hex[4:6])

	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func hexToByte(s string) uint8 {
	val, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return 0
	}
	return uint8(val)
}
