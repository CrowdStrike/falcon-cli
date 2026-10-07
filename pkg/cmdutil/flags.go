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

package cmdutil

import (
	"strings"

	"github.com/crowdstrike/falcon-cli/pkg/output"
	"github.com/spf13/cobra"
)

// AddFilterFlag adds a --filter flag for FQL expressions and registers a
// PreRunE that validates the filter syntax before the command runs.
func AddFilterFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "filter", "", "FQL filter expression (e.g. \"severity:'HIGH'\")")
	prev := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := ValidateFQLFilter(*p); err != nil {
			return err
		}
		if prev != nil {
			return prev(cmd, args)
		}
		return nil
	}
}

// AddPaginationFlags adds --limit and --offset flags.
func AddPaginationFlags(cmd *cobra.Command, limit *int64, offset *int64, defaultLimit int64) {
	cmd.Flags().Int64Var(limit, "limit", defaultLimit, "Maximum number of results to return")
	cmd.Flags().Int64Var(offset, "offset", 0, "Starting offset for results")
}

// AddSortFlag adds a --sort flag.
func AddSortFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "sort", "", "Sort expression (e.g. \"severity|desc\")")
}

// AddOutputFlag adds an --output flag (table|json|jsonl|toon).
func AddOutputFlag(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVar(p, "output", "table", "Output format: table, json, jsonl, or toon (compact, LLM-friendly)")
}

// AddIDsFlag adds an --ids flag for multi-ID commands.
func AddIDsFlag(cmd *cobra.Command, p *[]string, required bool) {
	cmd.Flags().StringSliceVar(p, "ids", nil, "One or more resource IDs (repeat flag or comma-separated)")
	if required {
		_ = cmd.MarkFlagRequired("ids")
	}
}

// TableFormattingOptions holds table rendering options populated by AddTableFormattingFlags.
type TableFormattingOptions struct {
	NoTruncate  bool
	Columns     string
	NoHeaders   bool
	MaxColWidth int
}

// AddTableFormattingFlags adds --no-truncate, --max-col-width, --columns, and --no-headers flags.
func AddTableFormattingFlags(cmd *cobra.Command, opts *TableFormattingOptions) {
	cmd.Flags().BoolVar(&opts.NoTruncate, "no-truncate", false, "Disable column value truncation")
	cmd.Flags().IntVar(&opts.MaxColWidth, "max-col-width", 50, "Maximum column width before truncation (0 = unlimited)")
	cmd.Flags().StringVar(&opts.Columns, "columns", "", "Comma-separated list of columns to display (e.g. id,name,provider)")
	cmd.Flags().BoolVar(&opts.NoHeaders, "no-headers", false, "Suppress table header row")
}

// ToOutputTableOptions converts TableFormattingOptions to output.TableOptions.
func (o *TableFormattingOptions) ToOutputTableOptions() *output.TableOptions {
	var cols []string
	if o.Columns != "" {
		for _, c := range strings.Split(o.Columns, ",") {
			c = strings.TrimSpace(c)
			if c != "" {
				cols = append(cols, c)
			}
		}
	}
	maxColWidth := o.MaxColWidth
	if maxColWidth == 0 {
		maxColWidth = 50
	}
	if o.NoTruncate {
		maxColWidth = 0
	}
	return &output.TableOptions{
		NoTruncate:  o.NoTruncate,
		Columns:     cols,
		NoHeaders:   o.NoHeaders,
		MaxColWidth: maxColWidth,
	}
}

// ClientFilterOptions holds client-side filter/sort options populated by AddClientFilterFlags.
type ClientFilterOptions struct {
	Grep   string
	SortBy string
}

// AddClientFilterFlags adds --grep and --sort-by flags for client-side filtering.
func AddClientFilterFlags(cmd *cobra.Command, opts *ClientFilterOptions) {
	cmd.Flags().StringVar(&opts.Grep, "grep", "", "Case-insensitive substring match across all columns")
	cmd.Flags().StringVar(&opts.SortBy, "sort-by", "", "Sort by column name (prefix with - for descending, e.g. -region)")
}

// ToFilterOptions converts ClientFilterOptions to output.FilterOptions.
func (o *ClientFilterOptions) ToFilterOptions() *output.FilterOptions {
	return &output.FilterOptions{
		Grep:   o.Grep,
		SortBy: o.SortBy,
	}
}
