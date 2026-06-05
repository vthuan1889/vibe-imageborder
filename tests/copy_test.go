package tests

import (
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	imgservice "vibe-imageborder/internal/image"
)

func TestCopyPipelinePreservesStructureAndChangesHash(t *testing.T) {
	source := t.TempDir()
	dest := t.TempDir()

	sub := filepath.Join(source, "album", "2024")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}

	srcPNG := filepath.Join(sub, "photo.png")
	srcJPG := filepath.Join(source, "cover.jpg")
	writePNG(srcPNG, color.RGBA{200, 100, 50, 255})
	writeMinimalJPG(t, srcJPG)

	svc := imgservice.NewService()
	files, err := imgservice.CollectImageFiles(source, 100)
	if err != nil {
		t.Fatalf("collect failed: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	for i, srcPath := range files {
		rel, err := filepath.Rel(source, srcPath)
		if err != nil {
			t.Fatal(err)
		}

		ext := filepath.Ext(srcPath)
		format, quality, outputExt := imgservice.FormatFromExtension(ext)
		destRel := rel[:len(rel)-len(ext)] + outputExt
		destPath := filepath.Join(dest, destRel)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			t.Fatal(err)
		}

		img, err := svc.LoadImage(srcPath)
		if err != nil {
			t.Fatal(err)
		}

		unique := imgservice.Uniquify(img, int64(i+1))
		outputBase := destPath[:len(destPath)-len(outputExt)]
		if err := svc.SaveImage(unique, outputBase, format, quality); err != nil {
			t.Fatal(err)
		}

		srcHash := fileMD5(t, srcPath)
		dstHash := fileMD5(t, destPath)
		if srcHash == dstHash {
			t.Errorf("hash should differ for %s", rel)
		}
	}

	if _, err := os.Stat(filepath.Join(dest, "album", "2024", "photo.png")); err != nil {
		t.Errorf("nested output missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "cover.jpg")); err != nil {
		t.Errorf("root output missing: %v", err)
	}
}

func writePNG(path string, c color.RGBA) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func writeMinimalJPG(t *testing.T, path string) {
	t.Helper()
	pngPath := path + ".tmp.png"
	writePNG(pngPath, color.RGBA{10, 20, 30, 255})
	defer os.Remove(pngPath)

	svc := imgservice.NewService()
	img, err := svc.LoadImage(pngPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.SaveImage(img, path[:len(path)-4], "jpg", 95); err != nil {
		t.Fatal(err)
	}
}

func fileMD5(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}
