package main

import (
	"fmt"
)

func humanSize(size int64) string {
	switch {
	case size >= GB:
		return fmt.Sprintf("%.1fG", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.1fM", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.1fK", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}

// calcBarsize determines the visual length of the bar.
// It subtracts the ratio's share from the total BarLength.
func calcBarsize(ratio float32) int {
	return BarLength - int(float32(BarLength)*ratio)
}

// Maps a ratio (0.0–1.0) to a color index (1–6) (Red to Cyan).
// Returns an integer from 1 to 6, ordered from highest to lowest.
func calcColor(ratio float32) int {
	color := int(ratio*6 + 1.01)
	return color
}

// calcRatio calculates the weight ratio (0.0 to 1.0).
// It returns the inverse progress: 1.0 - (currentSize / maxSize).
func calcRatio(size, max int64) float32 {
	return 1.0 - (float32(size) / float32(max))
}
