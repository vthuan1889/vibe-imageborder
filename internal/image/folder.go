package image

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// CollectImageFiles walks root recursively and returns absolute paths to image files.
func CollectImageFiles(root string, maxFiles int) ([]string, error) {
	root, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return nil, fmt.Errorf("invalid root path: %w", err)
	}

	var files []string
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !imageExtensions[ext] {
			return nil
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil
		}

		files = append(files, absPath)
		if len(files) > maxFiles {
			return fmt.Errorf("file count exceeds maximum %d", maxFiles)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no images found in source folder")
	}

	return files, nil
}

// FormatFromExtension returns save format and quality for a file extension.
// WebP sources are saved as PNG because the encoder has no native WebP support.
func FormatFromExtension(ext string) (format string, quality int, outputExt string) {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "jpg", 95, ext
	case ".png":
		return "png", 100, ext
	case ".webp":
		return "png", 100, ".png"
	default:
		return "png", 100, ".png"
	}
}
