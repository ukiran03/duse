package main

import (
	"strings"
	"testing"
)

func TestTableSort(t *testing.T) {
	tests := []struct {
		name  string
		rows  []*Row
		order int
		check func([]*Row) bool
	}{
		{
			name: "ascending sort",
			rows: []*Row{
				{size: 100, name: "large"},
				{size: 10, name: "small"},
				{size: 50, name: "medium"},
			},
			order: 1,
			check: func(rows []*Row) bool {
				return rows[0].size == 10 && rows[1].size == 50 &&
					rows[2].size == 100
			},
		},
		{
			name: "descending sort",
			rows: []*Row{
				{size: 10, name: "small"},
				{size: 100, name: "large"},
				{size: 50, name: "medium"},
			},
			order: -1,
			check: func(rows []*Row) bool {
				return rows[0].size == 100 && rows[1].size == 50 &&
					rows[2].size == 10
			},
		},
		{
			name: "no sort",
			rows: []*Row{
				{size: 100, name: "large"},
				{size: 10, name: "small"},
				{size: 50, name: "medium"},
			},
			order: 0,
			check: func(rows []*Row) bool {
				return rows[0].size == 100 && rows[1].size == 10 &&
					rows[2].size == 50
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{rows: tt.rows}
			table.SortTable(tt.order)
			if !tt.check(table.rows) {
				t.Errorf(
					"SortTable(%d) did not produce expected order",
					tt.order,
				)
			}
		})
	}
}

func TestTableSummarise(t *testing.T) {
	tests := []struct {
		name     string
		rows     []*Row
		expected int // expected number of rows after summarization
	}{
		{
			name: "only directories",
			rows: []*Row{
				{size: 100, name: "dir1/", isDir: true},
				{size: 200, name: "dir2/", isDir: true},
			},
			expected: 2,
		},
		{
			name: "only files",
			rows: []*Row{
				{size: 100, name: "file1", isDir: false},
				{size: 200, name: "file2", isDir: false},
			},
			expected: 1, // should be summarized to single "."
		},
		{
			name: "mixed",
			rows: []*Row{
				{size: 100, name: "dir1/", isDir: true},
				{size: 50, name: "file1", isDir: false},
				{size: 30, name: "file2", isDir: false},
				{size: 200, name: "dir2/", isDir: true},
			},
			expected: 3, // 2 dirs + 1 summarized file entry
		},
		{
			name:     "empty",
			rows:     []*Row{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{rows: tt.rows}
			table.SummariseTable()
			if len(table.rows) != tt.expected {
				t.Errorf(
					"SummariseTable() produced %d rows, want %d",
					len(table.rows),
					tt.expected,
				)
			}
		})
	}
}

func TestTableRecalculateBars(t *testing.T) {
	tests := []struct {
		name  string
		rows  []*Row
		check func([]*Row) bool
	}{
		{
			name: "normal case",
			rows: []*Row{
				{size: 100, name: "large"},
				{size: 10, name: "small"},
			},
			check: func(rows []*Row) bool {
				// After recalculation, bars should be set
				return rows[0].barlength >= 0 && rows[1].barlength >= 0 &&
					rows[0].color >= 1 && rows[1].color <= 6
			},
		},
		{
			name: "single row",
			rows: []*Row{
				{size: 100, name: "single"},
			},
			check: func(rows []*Row) bool {
				// Single row is both min and max,
				// so ratio = 0, barlength = 20, color = 1
				return rows[0].barlength == 20 && rows[0].color == 1
			},
		},
		{
			name:  "empty",
			rows:  []*Row{},
			check: func(rows []*Row) bool { return true },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{rows: tt.rows}
			table.recalculateBars()
			if !tt.check(table.rows) {
				t.Errorf("recalculateBars() did not produce expected results")
			}
		})
	}
}

func TestRowMakeBar(t *testing.T) {
	tests := []struct {
		name      string
		barlength int
		check     func(string) bool
	}{
		{
			name:      "full bar",
			barlength: 20,
			check: func(bar string) bool {
				// Should be all solid characters, no blanks
				return !strings.Contains(bar, BlankChar)
			},
		},
		{
			name:      "half bar",
			barlength: 10,
			check: func(bar string) bool {
				// Should have both blank and solid characters
				return strings.Contains(bar, BlankChar) &&
					strings.Contains(bar, SolidChar)
			},
		},
		{
			name:      "zero bar",
			barlength: 0,
			check: func(bar string) bool {
				// Should be all blank characters
				return !strings.Contains(bar, SolidChar) &&
					strings.Contains(bar, BlankChar)
			},
		},
		{
			name:      "overfull bar",
			barlength: 25,
			check: func(bar string) bool {
				// Should be all solid characters, no blanks
				return !strings.Contains(bar, BlankChar)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := &Row{barlength: tt.barlength}
			bar := row.makeBar()
			if !tt.check(bar) {
				t.Errorf(
					"makeBar() did not produce expected bar pattern for barlength %d",
					tt.barlength,
				)
			}
		})
	}
}
