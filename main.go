package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func main() {
	root := "/home/ukiran/Videos"
	depth := 1

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	for _, ent := range traverseFs(root, depth) {
		if ent.IsDir {
			fmt.Fprintf(w, "%s%s%s\t\t%s\n", Blue, humanSize(ent.Size), Reset, ent.Name)
		} else {
			fmt.Fprintf(w, "%s%s%s\t\t%s\n", Reset, humanSize(ent.Size), Reset, ent.Name)
		}
	}
	w.Flush()
}
