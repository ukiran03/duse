package main

const BarLength = 20

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
	White

	Text
	BoldText
	Reset
)

var colors = map[int]string{
	1: "\033[0031m",
	2: "\033[0032m",
	3: "\033[0033m",
	4: "\033[0034m",
	5: "\033[0035m",
	6: "\033[0036m",
	7: "\033[0037m",

	8:  "\033[0;39m", // Normal Text
	9:  "\033[1;39m", // Bold Text
	10: "\033[0000m",
}
