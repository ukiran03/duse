package main

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

const (
	Red int = iota + 1
	Green
	Yellow
	Blue
	Magenta
	Cyan

	Text
	Reset
)

var colors = map[int]string{
	1: "\033[031m", // Red
	2: "\033[032m", // Green
	3: "\033[033m", // Yellow
	4: "\033[034m", // Blue
	5: "\033[035m", // Magenta
	6: "\033[036m", // Cyan
	7: "\033[039m", // Text
	8: "\033[000m", // Reset
}
