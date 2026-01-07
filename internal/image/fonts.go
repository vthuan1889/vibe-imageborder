// Package image provides font loading and management.
package image

import (
	"embed"
	"fmt"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// FontManager handles font loading and caching.
type FontManager struct {
	fonts embed.FS
	cache map[string]*opentype.Font
	mu    sync.RWMutex
}

// NewFontManager creates font manager with embedded fonts.
func NewFontManager(fontsFS embed.FS) *FontManager {
	return &FontManager{
		fonts: fontsFS,
		cache: make(map[string]*opentype.Font),
	}
}

// LoadFont loads font from embedded FS.
func (fm *FontManager) LoadFont(name string) (*opentype.Font, error) {
	fm.mu.RLock()
	if cached, ok := fm.cache[name]; ok {
		fm.mu.RUnlock()
		return cached, nil
	}
	fm.mu.RUnlock()

	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Double-check after acquiring write lock
	if cached, ok := fm.cache[name]; ok {
		return cached, nil
	}

	path := "assets/fonts/" + name + ".ttf"
	data, err := fm.fonts.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read font %s: %w", name, err)
	}

	f, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font %s: %w", name, err)
	}

	fm.cache[name] = f
	return f, nil
}

// GetFace returns font.Face for given font and size.
func (fm *FontManager) GetFace(name string, size float64) (font.Face, error) {
	f, err := fm.LoadFont(name)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create face: %w", err)
	}

	return face, nil
}

// DefaultFontName returns default font to use.
func DefaultFontName() string {
	return "BeVietnamPro-Regular"
}
