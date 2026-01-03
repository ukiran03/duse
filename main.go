package main

import (
	"os"
)

func main() {
	root := "/home/ukiran/Git"
	depth := 1

	table := &Table{}
	table = writeBar(tabulate(traverseFs(root, depth)), BarLength)
	table.Print(os.Stdout)
}
