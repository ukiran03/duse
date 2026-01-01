package main

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const barLength = 20

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
		if row.isDir {
			fmt.Fprintf(
				w, "%s%s%s\t%s\t %s\n", Blue, humanSize(row.size),
				Reset, drawBar(row.barsize), row.name,
			)
		} else {
			fmt.Fprintf(
				w, "%s%s%s\t%s\t %s\n", Text, humanSize(row.size),
				Reset, drawBar(row.barsize), row.name,
			)
		}
	}
	w.Flush()
}

func humanSize(size int64) string {
	switch {
	case size >= GB:
		return fmt.Sprintf("%.1fG", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.1fM", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.1fK", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}

func drawBar(n int) string {
	blankBlock := " "
	charBlock := "▇"
	blanks := barLength - n
	return strings.Repeat(blankBlock, blanks) +
		strings.Repeat(charBlock, n)
}
