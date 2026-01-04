package main

import (
	"cmp"
	"fmt"
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

func (r *Row) makeBar() string {
	blankChar := " "
	solidChar := "▇"
	solidBar := strings.Repeat(solidChar, r.barlength)
	if r.barlength < BarLength {
		blankBar := strings.Repeat(blankChar, BarLength-r.barlength)
		return blankBar + solidBar
	}
	return solidBar
}

func (r *Row) humanSize() string {
	size := r.size
	return humanSize(size)
}

type Table struct {
	rows    []*Row
	largest int
}

func makeTable(entries []*FileEntry) *Table {
	if len(entries) == 0 {
		return &Table{}
	}
	dirs := make([]*Row, 0, len(entries)/2)
	files := make([]*Row, 0, len(entries)/2)
	for _, ent := range entries {
		r := &Row{
			name:  ent.Name,
			size:  ent.Size,
			isDir: ent.IsDir,
		}
		if ent.IsDir {
			dirs = append(dirs, r)
		} else {
			files = append(files, r)
		}
	}
	rows := slices.Concat(dirs, files)
	maxIdx := 0
	var maxSize int64
	for i, r := range rows {
		if r.size > maxSize {
			maxSize = r.size
			maxIdx = i
		}
	}
	for _, r := range rows {
		ratio := calcRatio(r.size, maxSize)
		r.barlength = calcBarsize(ratio)
		r.color = calcColor(ratio)
	}
	return &Table{rows, maxIdx}
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
	for _, row := range t.rows {
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

// Summarise merges all files into a single row while keeping
// directories separate.
// https://go.dev/wiki/SliceTricks#filtering-without-allocating
// https://abhinavg.net/2019/07/11/zero-alloc-slice-filter/
func (t *Table) SummariseTable() {
	if len(t.rows) == 0 {
		return
	}
	var summSize int64
	var hasFiles bool
	keeping := t.rows[:0]

	for _, row := range t.rows {
		if row.isDir {
			keeping = append(keeping, row)
		} else {
			summSize += row.size
			hasFiles = true
		}
	}
	if hasFiles {
		keeping = append(keeping, &Row{size: summSize, name: ".", isDir: false})
	}
	clear(t.rows[len(keeping):])
	t.rows = keeping

	var maxSize int64
	maxIdx := -1
	for i, row := range t.rows {
		if row.size >= maxSize {
			maxSize = row.size
			maxIdx = i
		}
	}
	for _, r := range t.rows {
		ratio := calcRatio(r.size, maxSize)
		r.barlength = calcBarsize(ratio)
		r.color = calcColor(ratio)
	}
	t.largest = maxIdx
}

func (t *Table) SortTable(order int) {
	if t == nil || len(t.rows) < 2 {
		return
	}

	if order > 0 {
		slices.SortFunc(t.rows, func(a, b *Row) int {
			if a == nil || b == nil {
				return 0
			}
			return cmp.Compare(a.size, b.size)
		})
	} else {
		slices.SortFunc(t.rows, func(b, a *Row) int {
			if a == nil || b == nil {
				return 0
			}
			return cmp.Compare(a.size, b.size)
		})
	}
}
