package main

func main() {
	root := "/home/ukiran/Videos" // test
	depth := 1

	table := makeTable(traverseFs(root, depth))
	// table.SortTable(-1)
	table.SummariseTable()
	table.Print()
}
