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
	flag.BoolVar(
		&cfg.summary,
		"m",
		false,
		"Summarise the individual files (.)",
	)
	flag.IntVar(&cfg.depth, "depth", 1, "Depth to traverse files")
	asc := flag.Bool("sort", false, "Sort in ascending order")
	desc := flag.Bool("rsort", false, "Sort in descending order")
	flag.Parse()

	paths := flag.Args()

	// Validate depth parameter
	if cfg.depth < 0 {
		fmt.Fprintf(os.Stderr, "Error: depth must be non-negative\n")
		os.Exit(1)
	}
	if len(paths) == 0 {
		paths = []string{"."}
	}
	cfg.paths = paths
	switch {
	case *asc && *desc:
		fmt.Fprintf(os.Stderr, "Error: cannot specify both -sort and -rsort\n")
		os.Exit(1)
	case *asc:
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
	var errors []error
	for _, dir := range cfg.paths {
		// basic validation check
		info, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				errors = append(
					errors,
					fmt.Errorf("path does not exist: %s", dir),
				)
			} else {
				errors = append(
					errors,
					fmt.Errorf("error accessing path %s: %v", dir, err),
				)
			}
			continue
		}
		if !info.IsDir() {
			errors = append(
				errors,
				fmt.Errorf("path is not a directory: %s", dir),
			)
			continue
		}

		table := makeTable(traverseFs(dir, cfg.depth))
		if cfg.summary {
			table.SummariseTable()
		}
		table.SortTable(cfg.sort)
		table.Print()
	}

	// Report all collected errors
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return fmt.Errorf("completed with %d error(s)", len(errors))
	}
	return nil
}
