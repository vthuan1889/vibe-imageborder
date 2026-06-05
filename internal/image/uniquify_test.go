package image

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestUniquifyPreservesDimensions(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 150))
	result := Uniquify(img, 42)

	w, h := NewService().GetDimensions(result)
	if w != 200 || h != 150 {
		t.Errorf("expected 200x150, got %dx%d", w, h)
	}
}

func TestUniquifyDifferentSeedsProduceDifferentOutput(t *testing.T) {
	img := solidImage(100, 100, color.RGBA{120, 80, 200, 255})

	out1 := encodePNG(Uniquify(img, 1))
	out2 := encodePNG(Uniquify(img, 2))

	if bytes.Equal(out1, out2) {
		t.Error("expected different PNG output for different seeds")
	}
}

func TestUniquifyMaxPixelDeltaIsOne(t *testing.T) {
	img := gradientImage(50, 50)
	result := Uniquify(img, 99)

	src := img.(*image.RGBA)
	changed := 0
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			srcRGBA := src.RGBAAt(x, y)
			dstRGBA := color.RGBAModel.Convert(result.At(x, y)).(color.RGBA)
			delta := MaxPixelDelta(srcRGBA, dstRGBA)
			if delta > 0 {
				changed++
			}
			if delta > 1 {
				t.Errorf("pixel (%d,%d) delta %d exceeds 1", x, y, delta)
			}
		}
	}
	if changed == 0 {
		t.Error("expected at least one pixel to change")
	}
}

func solidImage(w, h int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

func gradientImage(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8((x * 255) / w),
				G: uint8((y * 255) / h),
				B: uint8(((x + y) * 255) / (w + h)),
				A: 255,
			})
		}
	}
	return img
}

func encodePNG(img image.Image) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
