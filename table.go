package main

import (
	"cmp"
	"fmt"
	"iter"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
)

type Row struct {
	size      int64
	name      string
	isDir     bool
	barlength int
	color     int
}

type Table struct {
	rows []*Row
}

func (t *Table) All() iter.Seq[*Row] {
	return func(yield func(*Row) bool) {
		for _, row := range t.rows {
			if !yield(row) {
				return
			}
		}
	}
}

func EntriesByType(entries []*FileEntry) iter.Seq[*Row] {
	return func(yield func(*Row) bool) {
		// Pass 1: Yield Directories
		for _, ent := range entries {
			if ent.IsDir {
				if !yield(&Row{name: ent.Name, size: ent.Size, isDir: true}) {
					return
				}
			}
		}
		// Pass 2: Yield Files
		for _, ent := range entries {
			if !ent.IsDir {
				if !yield(&Row{name: ent.Name, size: ent.Size, isDir: false}) {
					return
				}
			}
		}
	}
}

func makeTable(entries []*FileEntry) *Table {
	if len(entries) == 0 {
		return &Table{}
	}

	// Drain the iterator straight into a single rows slice
	rows := slices.Collect(EntriesByType(entries))

	// Find max and calculate bars
	var maxSize int64
	for _, r := range rows {
		if r.size > maxSize {
			maxSize = r.size
		}
	}
	for _, r := range rows {
		ratio := calcRatio(r.size, maxSize)
		r.barlength = calcBarsize(ratio)
		r.color = calcColor(ratio)
	}
	return &Table{rows}
}

func (t *Table) SummariseTable() {
	if len(t.rows) == 0 {
		return
	}

	// Create an iterator that filters out files and aggregates their size
	summarizedSeq := func(yield func(*Row) bool) {
		var summSize int64
		var hasFiles bool

		for _, row := range t.rows {
			if row.isDir {
				if !yield(row) {
					return
				}
			} else {
				summSize += row.size
				hasFiles = true
			}
		}
		if hasFiles {
			yield(&Row{size: summSize, name: ".", isDir: false})
		}
	}

	// Collect the iterator results back into t.rows
	t.rows = slices.Collect(summarizedSeq)

	// Recalculate max sizes and bars
	var maxSize int64
	for _, row := range t.rows {
		if row.size >= maxSize {
			maxSize = row.size
		}
	}
	for _, r := range t.rows {
		ratio := calcRatio(r.size, maxSize)
		r.barlength = calcBarsize(ratio)
		r.color = calcColor(ratio)
	}
}

// SortTable sorts the table rows based on the provided order.
// Positive values sort ascending; negative values sort descending.
func (t *Table) SortTable(order int) {
	if t == nil || len(t.rows) < 2 {
		return
	}
	switch {
	case order > 0: // Ascending
		slices.SortFunc(t.rows, func(a, b *Row) int {
			return cmp.Compare(a.size, b.size)
		})
	case order < 0: // Descending
		slices.SortFunc(t.rows, func(a, b *Row) int {
			return cmp.Compare(b.size, a.size)
		})
	default:
		// No sort
		return
	}
}

func (r *Row) makeBar() string {
	solidBar := strings.Repeat(SolidChar, r.barlength)
	if r.barlength < BarLength {
		blankBar := strings.Repeat(BlankChar, BarLength-r.barlength)
		return blankBar + solidBar
	}
	return solidBar
}

func (r *Row) humanSize() string {
	size := r.size
	return humanSize(size)
}

func (t *Table) Total() string {
	var total int64
	for _, row := range t.rows {
		total += row.size
	}
	return humanSize(total)
}

func (t *Table) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)

	for row := range t.All() {
		fmt.Fprintf(
			w, "%s%s\t%s\t %s%s\n", colors[row.color], row.humanSize(),
			row.makeBar(), row.name, colors[Reset],
		)
	}
	fmt.Fprintf(
		w, "%s%s\t \t %s%s\n", colors[BoldText], t.Total(), "Total", colors[Reset],
	)
	w.Flush()
}
