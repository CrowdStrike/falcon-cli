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
	"testing"

	"github.com/go-openapi/runtime"
)

func TestToOutputTableOptions_ZeroWidthDisablesTruncation(t *testing.T) {
	opts := (&TableFormattingOptions{MaxColWidth: 0}).ToOutputTableOptions()
	if !opts.NoTruncate || opts.MaxColWidth != 0 {
		t.Errorf("got NoTruncate=%v MaxColWidth=%d, want true/0", opts.NoTruncate, opts.MaxColWidth)
	}

	opts = (&TableFormattingOptions{MaxColWidth: 20}).ToOutputTableOptions()
	if opts.NoTruncate || opts.MaxColWidth != 20 {
		t.Errorf("got NoTruncate=%v MaxColWidth=%d, want false/20", opts.NoTruncate, opts.MaxColWidth)
	}
}

func TestHandleAPIError_ScopeHintForCallerOperations(t *testing.T) {
	for _, op := range []string{"query IOM findings", "ListCloudGroupsExternal", "DeleteSuppressionRules"} {
		err := HandleAPIError(runtime.NewAPIError(op, nil, 403), op)
		if err == nil || !strings.Contains(err.Error(), "Required scope:") {
			t.Errorf("%s: expected scope hint, got %v", op, err)
		}
	}
}
