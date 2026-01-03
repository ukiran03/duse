package main

import (
	"fmt"
	"strings"
)

const BarLength = 20

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
	length := BarLength - int(BarLength*ratio)
	solidBar := strings.Repeat(charBlock, length)

	if length < BarLength {
		blankBar := strings.Repeat(blankBlock, BarLength-length)
		return fmt.Sprintf("%s%s%s", startC, blankBar+solidBar, endC)
	}
	return fmt.Sprintf("%s%s%s", startC, solidBar, endC)
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

func writeBar(t *Table, barLength int64) *Table {
	if len(t.rows) == 0 {
		return t
	}
	largest := t.rows[t.largest].size
	if largest == 0 {
		fmt.Println("Warning: Largest size is zero, skipping bar calculation.")
		return t
	}
	for _, r := range t.rows {
		cols := (BarLength * r.size) / largest
		r.barlength = int(cols)
	}
	return t
}
