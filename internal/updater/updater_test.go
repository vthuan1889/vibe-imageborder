package updater

import "testing"

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"equal versions", "v1.0.0", "v1.0.0", 0},
		{"v1 less than v2 major", "v1.0.0", "v2.0.0", -1},
		{"v1 greater than v2 major", "v2.0.0", "v1.0.0", 1},
		{"v1 less than v2 minor", "v1.0.0", "v1.1.0", -1},
		{"v1 greater than v2 minor", "v1.2.0", "v1.1.0", 1},
		{"v1 less than v2 patch", "v1.0.0", "v1.0.1", -1},
		{"v1 greater than v2 patch", "v1.0.2", "v1.0.1", 1},
		{"without v prefix", "1.0.0", "1.0.1", -1},
		{"mixed prefix", "v1.0.0", "1.0.1", -1},
		{"dev version", "dev", "v1.0.0", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("CompareVersions(%q, %q) = %d; want %d",
					tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}
