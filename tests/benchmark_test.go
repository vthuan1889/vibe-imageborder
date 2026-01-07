package tests

import (
	"image/color"
	"testing"

	imgservice "vibe-imageborder/internal/image"
)

func BenchmarkComposite(b *testing.B) {
	// Create small test images for benchmark
	imageSvc := imgservice.NewService()
	compositor := imgservice.NewCompositor(imageSvc)

	// Create test images in memory
	product := createTestImage(800, 800, color.RGBA{255, 0, 0, 255})
	frame := createTestImage(1000, 1000, color.RGBA{0, 0, 255, 128})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compositor.Composite(product, frame, "#ffffff")
	}
}

func BenchmarkResizeToFit(b *testing.B) {
	imageSvc := imgservice.NewService()
	product := createTestImage(2000, 2000, color.RGBA{255, 0, 0, 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		imageSvc.ResizeToFit(product, 1000, 1000)
	}
}
