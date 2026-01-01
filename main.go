package main

import (
	"os"
)

func main() {
	root := "/home/ukiran/Videos/Tuts"
	depth := 1

	table := &Table{}
	table = writeBar(tabulate(traverseFs(root, depth)), barLength)
	table.Print(os.Stdout)
}
