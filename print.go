package main

import (
	"fmt"
	"strings"
)

const barLength = 20

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

func drawBar(ratio float32) string {
	blankBlock := " "
	charBlock := "▇"
	startC := colors[calcColor(ratio)]
	endC := colors[Reset]
	length := barLength - int(barLength*ratio)
	solidBar := strings.Repeat(charBlock, length)

	if length < barLength {
		blankBar := strings.Repeat(blankBlock, barLength-length)
		return fmt.Sprintf("%s%s%s", startC, blankBar+solidBar, endC)
	}
	return fmt.Sprintf("%s%s%s", startC, solidBar, endC)
}

func calcColor(ratio float32) int {
	color := int(ratio*6 + 1.01)
	return color
}

func calcRatio(size, max int64) float32 {
	return 1.0 - (float32(size) / float32(max))
}
