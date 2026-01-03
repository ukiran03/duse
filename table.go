package main

import (
	"fmt"
	"io"
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
	rows    []*Row
	largest int
}

func makeTable(entries []*FileEntry) *Table {
	if len(entries) == 0 {
		return &Table{}
	}
	rows := make([]*Row, len(entries))
	var maxEntry *FileEntry
	maxIdx := 0
	for i, ent := range entries {
		rows[i] = &Row{
			name:  ent.Name,
			size:  ent.Size,
			isDir: ent.IsDir,
		}
		if maxEntry == nil || ent.Size > maxEntry.Size {
			maxEntry = ent
			maxIdx = i
		}
	}
	maxSize := maxEntry.Size
	for _, r := range rows {
		ratio := calcRatio(r.size, maxSize)
		r.barlength = calcBarsize(ratio)
		r.color = calcColor(ratio)
	}
	return &Table{rows, maxIdx}
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
