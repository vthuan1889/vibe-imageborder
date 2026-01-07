package image

import (
	"image/color"
	"testing"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		input    string
		expected color.Color
	}{
		{"white", color.RGBA{255, 255, 255, 255}},
		{"black", color.RGBA{0, 0, 0, 255}},
		{"red", color.RGBA{255, 0, 0, 255}},
		{"#FFFFFF", color.RGBA{255, 255, 255, 255}},
		{"#FF0000", color.RGBA{255, 0, 0, 255}},
		{"#FF0000AA", color.RGBA{255, 0, 0, 170}},
		{"#FFF", color.RGBA{255, 255, 255, 255}}, // Short hex
		{"invalid", color.RGBA{255, 255, 255, 255}}, // Default white
	}

	for _, tt := range tests {
		result := parseColor(tt.input)
		if result != tt.expected {
			t.Errorf("parseColor(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestParsePosition(t *testing.T) {
	tests := []struct {
		input     string
		expectX   int
		expectY   int
		expectErr bool
	}{
		{"98,1720", 98, 1720, false},
		{"0,0", 0, 0, false},
		{"100, 200", 100, 200, false}, // With spaces
		{"invalid", 0, 0, true},
		{"100", 0, 0, true}, // Missing y
	}

	for _, tt := range tests {
		x, y, err := parsePosition(tt.input)
		if tt.expectErr {
			if err == nil {
				t.Errorf("parsePosition(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("parsePosition(%q) unexpected error: %v", tt.input, err)
			}
			if x != tt.expectX || y != tt.expectY {
				t.Errorf("parsePosition(%q) = (%d, %d), want (%d, %d)",
					tt.input, x, y, tt.expectX, tt.expectY)
			}
		}
	}
}

func TestParseFontSize(t *testing.T) {
	tests := []struct {
		input     string
		expected  float64
		expectErr bool
	}{
		{"45", 45.0, false},
		{"12.5", 12.5, false},
		{" 30 ", 30.0, false}, // With spaces
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		size, err := parseFontSize(tt.input)
		if tt.expectErr {
			if err == nil {
				t.Errorf("parseFontSize(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("parseFontSize(%q) unexpected error: %v", tt.input, err)
			}
			if size != tt.expected {
				t.Errorf("parseFontSize(%q) = %f, want %f", tt.input, size, tt.expected)
			}
		}
	}
}

func TestReplaceVariables(t *testing.T) {
	tests := []struct {
		input    string
		values   map[string]string
		expected string
	}{
		{
			"[barcode]",
			map[string]string{"barcode": "ABC123"},
			"ABC123",
		},
		{
			"D[size_dai] x R[size_rong] x C[size_cao] cm",
			map[string]string{
				"size_dai":  "30",
				"size_rong": "20",
				"size_cao":  "15",
			},
			"D30 x R20 x C15 cm",
		},
		{
			"[unknown]",
			map[string]string{"barcode": "ABC123"},
			"[unknown]", // Unchanged if not in map
		},
		{
			"No variables here",
			map[string]string{},
			"No variables here",
		},
	}

	for _, tt := range tests {
		result := replaceVariables(tt.input, tt.values)
		if result != tt.expected {
			t.Errorf("replaceVariables(%q, %v) = %q, want %q",
				tt.input, tt.values, result, tt.expected)
		}
	}
}
