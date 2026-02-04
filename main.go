package main

import (
	"flag"
	"fmt"
	"os"
)

type config struct {
	summary bool
	sort    int
	depth   int
	paths   []string
}

func main() {
	cfg := config{}
	flag.BoolVar(&cfg.summary, "s", false, "Summarise the individual files (.)")
	flag.IntVar(&cfg.depth, "depth", 1, "Depth to traverse files")
	asc := flag.Bool("sort", false, "Sort in ascending order")
	desc := flag.Bool("rsort", false, "Sort in descending order")
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}
	cfg.paths = paths
	switch {
	case *asc && *desc || *asc:
		cfg.sort = 1
	case *desc:
		cfg.sort = -1
	default:
		cfg.sort = 0 // just in case (no sort)
	}

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1) // TODO: it must continue for other dirs
	}
}

func run(cfg config) error {
	for _, dir := range cfg.paths {
		// basic validation check
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dir)
		}

		table := makeTable(traverseFs(dir, cfg.depth))
		if cfg.summary {
			table.SummariseTable()
		}
		table.SortTable(cfg.sort)
		table.Print()
	}
	return nil // TODO: Improve error propagation
}
