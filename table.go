package main

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type Row struct {
	size    int64
	name    string
	isDir   bool
	barsize int
}

type Table struct {
	rows    []*Row
	largest int
}

func tabulate(ents []*FileEntry) *Table {
	n := len(ents)
	if n == 0 {
		return &Table{}
	}
	rows := make([]*Row, n)
	largest := 0
	for i, ent := range ents {
		if ent.Size > ents[largest].Size {
			largest = i
		}
		rows[i] = &Row{
			name:  ent.Name,
			size:  ent.Size,
			isDir: ent.IsDir,
		}
	}
	return &Table{rows, largest}
}

func writeBar(t *Table, barLength int64) *Table {
	if len(t.rows) == 0 {
		return t
	}
	largest := t.rows[t.largest].size
	if largest == 0 {
		fmt.Println("Warning: Largest size is zero, skipping bar calculation.")
		return t
	}
	for _, r := range t.rows {
		cols := (barLength * r.size) / largest
		r.barsize = int(cols)
	}
	return t
}

func (t *Table) Print(out io.Writer) {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', tabwriter.AlignRight)
	for _, row := range t.rows {
		ratio := calcRatio(row.size, t.rows[t.largest].size)
		if row.isDir {
			fmt.Fprintf(
				w, "%s%s%s\t%s\t %s%s%s\n", colors[Blue], humanSize(row.size),
				colors[Reset], drawBar(ratio), colors[Text], row.name, colors[Reset],
			)
		} else {
			fmt.Fprintf(
				w, "%s%s%s\t%s\t %s%s%s\n", colors[Text], humanSize(row.size),
				colors[Reset], drawBar(ratio), colors[Text], row.name, colors[Reset],
			)
		}
	}
	w.Flush()
}

// Summarise merges all files into a single row while keeping
// directories separate.
func (t *Table) SummariseInPlace() {
	if len(t.rows) == 0 {
		return
	}
	var sumSize int64
	var sumBar int
	var hasFiles bool

	// writeIdx keeps track of where to put the next directory
	writeIdx := 0
	for i := 0; i < len(t.rows); i++ {
		if t.rows[i].isDir {
			// Move directory to the "keep" section
			t.rows[writeIdx] = t.rows[i]
			writeIdx++
		} else {
			sumSize += t.rows[i].size
			sumBar += t.rows[i].barsize
			hasFiles = true
		}
	}
	if hasFiles {
		summaryRow := &Row{sumSize, ".", false, sumBar}
		// If there's space in the original capacity, we append/assign
		if writeIdx < len(t.rows) {
			t.rows[writeIdx] = summaryRow
			writeIdx++
		} else {
			t.rows = append(t.rows, summaryRow)
		}
	}

	// Empty remaining: nil out, resilce
	for i := writeIdx; i < len(t.rows); i++ {
		t.rows[i] = nil
	}
	t.rows = t.rows[:writeIdx]
}
