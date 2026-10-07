// Copyright (c) 2026 CrowdStrike, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package output

import (
	"io"
	"strings"
)

// Format represents the output format.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatJSONL Format = "jsonl"
	FormatTOON  Format = "toon"
)

// TableDefinition describes how to render a slice of items as a table.
type TableDefinition struct {
	// Headers is the list of column headers.
	Headers []string
	// RowFunc extracts column values from a single item.
	RowFunc func(item any) []string
}

// TableOptions controls table rendering behaviour.
type TableOptions struct {
	// NoTruncate disables column truncation (default: truncate to MaxColWidth).
	NoTruncate bool
	// MaxColWidth is the maximum character width per column (0 = unlimited / no truncation).
	MaxColWidth int
	// Columns, when non-empty, restricts output to the named columns (case-insensitive).
	Columns []string
	// NoHeaders suppresses the header row.
	NoHeaders bool
}

// FilterOptions controls client-side filtering and sorting applied after fetching.
type FilterOptions struct {
	// Grep filters rows to those where any column value contains this substring (case-insensitive).
	Grep string
	// SortBy sorts rows by the named column (case-insensitive). Prefix with "-" for descending.
	SortBy string
}

// Printer writes structured data to a writer.
type Printer interface {
	Print(w io.Writer, data any) error
}

// NewPrinter returns a Printer for the given format.
func NewPrinter(format Format, def *TableDefinition, tableOpts ...*TableOptions) Printer {
	var opts *TableOptions
	if len(tableOpts) > 0 && tableOpts[0] != nil {
		opts = tableOpts[0]
	} else {
		opts = &TableOptions{}
	}
	if opts.MaxColWidth == 0 && !opts.NoTruncate {
		opts.MaxColWidth = 50
	}

	switch format {
	case FormatJSON:
		return &jsonPrinter{}
	case FormatJSONL:
		return &jsonlPrinter{}
	case FormatTOON:
		return &toonPrinter{}
	default:
		return &tablePrinter{def: def, opts: opts}
	}
}

// columnIndices returns the subset of header indices that match the requested columns.
// If columns is empty all indices are returned.
func columnIndices(headers []string, columns []string) []int {
	if len(columns) == 0 {
		indices := make([]int, len(headers))
		for i := range headers {
			indices[i] = i
		}
		return indices
	}
	var indices []int
	for _, col := range columns {
		col = strings.ToLower(strings.TrimSpace(col))
		for i, h := range headers {
			if strings.ToLower(h) == col {
				indices = append(indices, i)
				break
			}
		}
	}
	return indices
}

// truncate returns s truncated to maxLen with "…" suffix if needed.
func truncate(s string, maxLen int) string {
	if maxLen <= 0 || len([]rune(s)) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen-1]) + "…"
}
