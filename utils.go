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

func calcBarsize(ratio float32) int {
	return BarLength - int(BarLength*ratio)
}

func calcColor(ratio float32) int {
	color := int(ratio*6 + 1.01)
	return color
}

func calcRatio(size, max int64) float32 {
	return 1.0 - (float32(size) / float32(max))
}
