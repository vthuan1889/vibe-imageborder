package image

import (
	"image/color"
	"testing"
)

func TestParsePosition(t *testing.T) {
	tests := []struct {
		input string
		x, y  int
		err   bool
	}{
		{"100,200", 100, 200, false},
		{"0,0", 0, 0, false},
		{" 50 , 100 ", 50, 100, false},
		{"invalid", 0, 0, true},
		{"100", 0, 0, true},
		{"a,b", 0, 0, true},
		{"", 0, 0, true},
	}

	for _, tt := range tests {
		x, y, err := ParsePosition(tt.input)
		if tt.err && err == nil {
			t.Errorf("Expected error for %s", tt.input)
		}
		if !tt.err && err != nil {
			t.Errorf("Unexpected error for %s: %v", tt.input, err)
		}
		if !tt.err && (x != tt.x || y != tt.y) {
			t.Errorf("For %s: expected %d,%d got %d,%d", tt.input, tt.x, tt.y, x, y)
		}
	}
}

func TestParseColorName(t *testing.T) {
	tests := []struct {
		input    string
		expected color.RGBA
	}{
		{"white", color.RGBA{255, 255, 255, 255}},
		{"WHITE", color.RGBA{255, 255, 255, 255}},
		{"black", color.RGBA{0, 0, 0, 255}},
		{"red", color.RGBA{255, 0, 0, 255}},
		{"green", color.RGBA{0, 255, 0, 255}},
		{"blue", color.RGBA{0, 0, 255, 255}},
		{"#ff0000", color.RGBA{255, 0, 0, 255}},
		{"#00ff00", color.RGBA{0, 255, 0, 255}},
		{"#0000ff", color.RGBA{0, 0, 255, 255}},
		{"unknown", color.RGBA{255, 255, 255, 255}}, // fallback to white
	}

	for _, tt := range tests {
		result := ParseColorName(tt.input)
		r1, g1, b1, a1 := result.RGBA()
		r2, g2, b2, a2 := tt.expected.RGBA()
		if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
			t.Errorf("For %s: expected %v got %v", tt.input, tt.expected, result)
		}
	}
}

func TestParseColorNameHex(t *testing.T) {
	tests := []struct {
		hex      string
		expected color.RGBA
	}{
		{"#ffffff", color.RGBA{255, 255, 255, 255}},
		{"#000000", color.RGBA{0, 0, 0, 255}},
		{"#f1eeea", color.RGBA{241, 238, 234, 255}},
	}

	for _, tt := range tests {
		result := ParseColorName(tt.hex)
		if rgba, ok := result.(color.RGBA); ok {
			if rgba != tt.expected {
				t.Errorf("For %s: expected %v got %v", tt.hex, tt.expected, rgba)
			}
		}
	}
}

func TestNewTextRenderer(t *testing.T) {
	// Test that TextRenderer can be created without FontManager
	// (will fail on DrawOverlays without proper fonts)
	tr := NewTextRenderer(nil)
	if tr == nil {
		t.Error("NewTextRenderer returned nil")
	}
}

func TestDefaultFontName(t *testing.T) {
	name := DefaultFontName()
	if name != "BeVietnamPro-Regular" {
		t.Errorf("Expected BeVietnamPro-Regular, got %s", name)
	}
}
