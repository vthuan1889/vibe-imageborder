package image

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectImageFilesRecursive(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "photos", "a")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}

	files := map[string][]byte{
		filepath.Join(root, "top.jpg"):       {0},
		filepath.Join(sub, "nested.png"):     {0},
		filepath.Join(root, "skip.txt"):      {0},
		filepath.Join(root, "hidden", ".x"):  {0},
	}
	_ = os.MkdirAll(filepath.Join(root, "hidden"), 0755)

	for path, data := range files {
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatal(err)
		}
	}

	collected, err := CollectImageFiles(root, 100)
	if err != nil {
		t.Fatalf("CollectImageFiles failed: %v", err)
	}

	if len(collected) != 2 {
		t.Fatalf("expected 2 images, got %d: %v", len(collected), collected)
	}
}

func TestFormatFromExtension(t *testing.T) {
	tests := []struct {
		ext        string
		format     string
		outputExt  string
	}{
		{".jpg", "jpg", ".jpg"},
		{".jpeg", "jpg", ".jpeg"},
		{".png", "png", ".png"},
		{".webp", "png", ".png"},
	}

	for _, tc := range tests {
		format, _, outputExt := FormatFromExtension(tc.ext)
		if format != tc.format || outputExt != tc.outputExt {
			t.Errorf("ext %s: got format=%s outputExt=%s, want format=%s outputExt=%s",
				tc.ext, format, outputExt, tc.format, tc.outputExt)
		}
	}
}
