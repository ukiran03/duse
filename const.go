package main

const (
	BarLength = 20
	BlankChar = " "
	SolidChar = "▇"
)

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
	Red:     "\x1b[31m",
	Green:   "\x1b[32m",
	Yellow:  "\x1b[33m",
	Blue:    "\x1b[34m",
	Magenta: "\x1b[35m",
	Cyan:    "\x1b[36m",
	White:   "\x1b[37m",

	Text:     "\x1b[39m", // Default foreground color
	BoldText: "\x1b[1m",  // Bold modifier
	Reset:    "\x1b[0m",  // Reset all attributes
}
