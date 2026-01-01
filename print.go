package main

import "fmt"

const (
	KB        = 1024
	MB        = KB * 1024
	GB        = MB * 1024
	barLength = 20
)

// ANSI color codes
const (
	Black = "\033[30m"
	Blue  = "\033[34m"
	Green = "\033[32m"
	Reset = "\033[0m"
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
