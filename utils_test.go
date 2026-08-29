package main

import (
	"testing"
)

func TestHumanSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{"bytes", 500, "500B"},
		{"kilobytes", 1536, "1.5K"},
		{"megabytes", 1572864, "1.5M"},
		{"gigabytes", 1610612736, "1.5G"},
		{"zero", 0, "0B"},
		{"exact KB", 1024, "1.0K"},
		{"exact MB", 1048576, "1.0M"},
		{"exact GB", 1073741824, "1.0G"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := humanSize(tt.size)
			if result != tt.expected {
				t.Errorf(
					"humanSize(%d) = %s, want %s",
					tt.size,
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestCalcBarsize(t *testing.T) {
	tests := []struct {
		name     string
		ratio    float32
		expected int
	}{
		{"zero ratio", 0.0, 20},
		{"half ratio", 0.5, 10},
		{"full ratio", 1.0, 0},
		{"small ratio", 0.1, 18},
		{"large ratio", 0.9, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calcBarsize(tt.ratio)
			if result != tt.expected {
				t.Errorf(
					"calcBarsize(%f) = %d, want %d",
					tt.ratio,
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestCalcColor(t *testing.T) {
	tests := []struct {
		name  string
		ratio float32
		min   int
		max   int
	}{
		{"zero ratio", 0.0, 1, 1},
		{"half ratio", 0.5, 4, 4},
		{"full ratio", 1.0, 7, 7},
		{"small ratio", 0.1, 1, 1},
		{"large ratio", 0.9, 6, 6},
		{"negative ratio", -0.1, 0, 0},
		{"ratio above 1", 1.1, 7, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calcColor(tt.ratio)
			if result < tt.min || result > tt.max {
				t.Errorf(
					"calcColor(%f) = %d, want between %d and %d",
					tt.ratio,
					result,
					tt.min,
					tt.max,
				)
			}
		})
	}
}

func TestCalcRatio(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		max      int64
		expected float32
	}{
		{"zero size", 0, 100, 1.0},
		{"max size", 100, 100, 0.0},
		{"half size", 50, 100, 0.5},
		{"small size", 10, 100, 0.9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calcRatio(tt.size, tt.max)
			if result != tt.expected {
				t.Errorf(
					"calcRatio(%d, %d) = %f, want %f",
					tt.size,
					tt.max,
					result,
					tt.expected,
				)
			}
		})
	}
}
