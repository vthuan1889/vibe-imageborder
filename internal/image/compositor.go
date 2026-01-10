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
// Uses the smaller dimensions between product and frame as target size.
func (c *Compositor) Composite(product, frame image.Image, bgColor string) *CompositeResult {
	productBounds := product.Bounds()
	frameBounds := frame.Bounds()

	// Use smaller dimensions between product and frame
	width := min(productBounds.Dx(), frameBounds.Dx())
	height := min(productBounds.Dy(), frameBounds.Dy())

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

	// Apply background color if specified
	if bgColor != "" {
		bg := c.service.CreateBlankCanvas(width, height, bgColor)
		draw.Draw(canvas, canvas.Bounds(), bg, image.Point{}, draw.Src)
	}

	// Resize product to fit target dimensions
	resizedProduct := c.service.ResizeToFit(product, width, height)

	// Resize frame to fit target dimensions
	resizedFrame := c.service.ResizeToFit(frame, width, height)

	// Center product on canvas
	productRect := resizedProduct.Bounds()
	offsetX := (width - productRect.Dx()) / 2
	offsetY := (height - productRect.Dy()) / 2

	// Draw product first (background)
	draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+productRect.Dx(), offsetY+productRect.Dy()),
		resizedProduct, image.Point{}, draw.Over)

	// Center frame on canvas
	frameRect := resizedFrame.Bounds()
	frameOffsetX := (width - frameRect.Dx()) / 2
	frameOffsetY := (height - frameRect.Dy()) / 2

	// Draw frame overlay (with alpha)
	draw.Draw(canvas, image.Rect(frameOffsetX, frameOffsetY, frameOffsetX+frameRect.Dx(), frameOffsetY+frameRect.Dy()),
		resizedFrame, image.Point{}, draw.Over)

	return &CompositeResult{
		Image:  canvas,
		Width:  width,
		Height: height,
	}
}

// CompositeWithPosition allows custom product positioning.
func (c *Compositor) CompositeWithPosition(product, frame image.Image, bgColor string, position string) *CompositeResult {
	productBounds := product.Bounds()
	frameBounds := frame.Bounds()

	// Use smaller dimensions between product and frame
	width := min(productBounds.Dx(), frameBounds.Dx())
	height := min(productBounds.Dy(), frameBounds.Dy())

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

	// Apply background if specified
	if bgColor != "" {
		bg := c.service.CreateBlankCanvas(width, height, bgColor)
		draw.Draw(canvas, canvas.Bounds(), bg, image.Point{}, draw.Src)
	}

	// Resize both images to fit target dimensions
	resizedProduct := c.service.ResizeToFit(product, width, height)
	resizedFrame := c.service.ResizeToFit(frame, width, height)

	// Calculate center positions
	productRect := resizedProduct.Bounds()
	offsetX := (width - productRect.Dx()) / 2
	offsetY := (height - productRect.Dy()) / 2

	frameRect := resizedFrame.Bounds()
	frameOffsetX := (width - frameRect.Dx()) / 2
	frameOffsetY := (height - frameRect.Dy()) / 2

	// Handle position-based ordering
	if position == "below" {
		// Frame first, then product
		draw.Draw(canvas, image.Rect(frameOffsetX, frameOffsetY, frameOffsetX+frameRect.Dx(), frameOffsetY+frameRect.Dy()),
			resizedFrame, image.Point{}, draw.Over)
		draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+productRect.Dx(), offsetY+productRect.Dy()),
			resizedProduct, image.Point{}, draw.Over)
	} else {
		// Product first, then frame (default)
		draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+productRect.Dx(), offsetY+productRect.Dy()),
			resizedProduct, image.Point{}, draw.Over)
		draw.Draw(canvas, image.Rect(frameOffsetX, frameOffsetY, frameOffsetX+frameRect.Dx(), frameOffsetY+frameRect.Dy()),
			resizedFrame, image.Point{}, draw.Over)
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
