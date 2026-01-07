// Package image provides image compositing operations.
package image

import (
	"fmt"
	"image"
	"image/draw"

	"vibe-imageborder/internal/models"
)

// Compositor handles image compositing.
type Compositor struct {
	service *Service
}

// NewCompositor creates new compositor.
func NewCompositor(service *Service) *Compositor {
	return &Compositor{service: service}
}

// CompositeResult holds the composited image.
type CompositeResult struct {
	Image  image.Image
	Width  int
	Height int
}

// Composite combines product and frame images.
// Product is resized to fit frame, then frame overlaid on top.
func (c *Compositor) Composite(product, frame image.Image, bgColor string) *CompositeResult {
	frameBounds := frame.Bounds()
	width := frameBounds.Dx()
	height := frameBounds.Dy()

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

	// Apply background color if specified
	if bgColor != "" {
		bg := c.service.CreateBlankCanvas(width, height, bgColor)
		draw.Draw(canvas, canvas.Bounds(), bg, image.Point{}, draw.Src)
	}

	// Resize product to fit frame dimensions
	resizedProduct := c.service.ResizeToFit(product, width, height)

	// Center product on canvas
	productBounds := resizedProduct.Bounds()
	offsetX := (width - productBounds.Dx()) / 2
	offsetY := (height - productBounds.Dy()) / 2

	// Draw product first (background)
	draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+productBounds.Dx(), offsetY+productBounds.Dy()),
		resizedProduct, image.Point{}, draw.Over)

	// Draw frame overlay (with alpha)
	draw.Draw(canvas, canvas.Bounds(), frame, image.Point{}, draw.Over)

	return &CompositeResult{
		Image:  canvas,
		Width:  width,
		Height: height,
	}
}

// CompositeWithPosition allows custom product positioning.
func (c *Compositor) CompositeWithPosition(product, frame image.Image, bgColor string, position string) *CompositeResult {
	frameBounds := frame.Bounds()
	width := frameBounds.Dx()
	height := frameBounds.Dy()

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

	// Apply background if specified
	if bgColor != "" {
		bg := c.service.CreateBlankCanvas(width, height, bgColor)
		draw.Draw(canvas, canvas.Bounds(), bg, image.Point{}, draw.Src)
	}

	// Resize product to fit
	resizedProduct := c.service.ResizeToFit(product, width, height)
	productBounds := resizedProduct.Bounds()

	// Center product
	offsetX := (width - productBounds.Dx()) / 2
	offsetY := (height - productBounds.Dy()) / 2

	// Handle position-based ordering
	if position == "below" {
		// Frame first, then product
		draw.Draw(canvas, canvas.Bounds(), frame, image.Point{}, draw.Over)
		draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+productBounds.Dx(), offsetY+productBounds.Dy()),
			resizedProduct, image.Point{}, draw.Over)
	} else {
		// Product first, then frame (default)
		draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+productBounds.Dx(), offsetY+productBounds.Dy()),
			resizedProduct, image.Point{}, draw.Over)
		draw.Draw(canvas, canvas.Bounds(), frame, image.Point{}, draw.Over)
	}

	return &CompositeResult{
		Image:  canvas,
		Width:  width,
		Height: height,
	}
}

// ToRGBA converts image to RGBA for drawing operations.
func ToRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}

// CompositeWithText combines product, frame, and text overlays.
func (c *Compositor) CompositeWithText(
	product, frame image.Image,
	bgColor string,
	overlays map[string]models.TextOverlay,
	textRenderer *TextRenderer,
) (*CompositeResult, error) {
	// First composite product + frame
	result := c.Composite(product, frame, bgColor)

	// Then draw text overlays
	if textRenderer != nil && len(overlays) > 0 {
		imgWithText, err := textRenderer.DrawOverlays(result.Image, overlays)
		if err != nil {
			return nil, fmt.Errorf("failed to draw text: %w", err)
		}
		result.Image = imgWithText
	}

	return result, nil
}
