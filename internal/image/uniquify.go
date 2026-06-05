package image

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/disintegration/imaging"
)

const (
	minUniquifyPixels = 50
	uniquifyRatio     = 10000 // ~0.01% of pixels
)

// Uniquify applies imperceptible LSB noise to change image fingerprint
// while keeping visual appearance identical to the human eye.
func Uniquify(img image.Image, seed int64) image.Image {
	nrgba := imaging.Clone(img)
	bounds := nrgba.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	numPixels := width * height / uniquifyRatio
	if numPixels < minUniquifyPixels {
		numPixels = minUniquifyPixels
	}

	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < numPixels; i++ {
		x := bounds.Min.X + rng.Intn(width)
		y := bounds.Min.Y + rng.Intn(height)
		c := nrgba.NRGBAAt(x, y)

		channel := rng.Intn(3)
		delta := int8(1)
		if rng.Intn(2) == 0 {
			delta = -1
		}

		switch channel {
		case 0:
			c.R = clampUint8(int(c.R) + int(delta))
		case 1:
			c.G = clampUint8(int(c.G) + int(delta))
		case 2:
			c.B = clampUint8(int(c.B) + int(delta))
		}

		nrgba.SetNRGBA(x, y, c)
	}

	return nrgba
}

func clampUint8(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

// MaxPixelDelta returns the maximum absolute difference between two colors.
func MaxPixelDelta(a, b color.RGBA) int {
	max := 0
	for _, d := range []int{int(a.R) - int(b.R), int(a.G) - int(b.G), int(a.B) - int(b.B)} {
		if d < 0 {
			d = -d
		}
		if d > max {
			max = d
		}
	}
	return max
}
