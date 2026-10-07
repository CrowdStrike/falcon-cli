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
	"bytes"
	"fmt"
	"testing"
	"unicode"
)

type toonItem struct {
	ID       string   `json:"id"`
	Severity string   `json:"severity"`
	Count    int64    `json:"count"`
	Tags     []string `json:"tags,omitempty"`
}

func TestToonPrinter_TabularArray(t *testing.T) {
	items := []toonItem{
		{ID: "a1", Severity: "HIGH", Count: 3},
		{ID: "b2", Severity: "LOW", Count: 10},
	}

	var buf bytes.Buffer
	if err := NewPrinter(FormatTOON, nil).Print(&buf, items); err != nil {
		t.Fatalf("Print: %v", err)
	}

	want := "[2]{id,severity,count}:\n  a1,HIGH,3\n  b2,LOW,10\n"
	if got := buf.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestToonPrinter_ObjectKeepsFieldOrder(t *testing.T) {
	var buf bytes.Buffer
	item := toonItem{ID: "a1", Severity: "HIGH", Count: 3, Tags: []string{"x", "y"}}
	if err := NewPrinter(FormatTOON, nil).Print(&buf, item); err != nil {
		t.Fatalf("Print: %v", err)
	}

	want := "id: a1\nseverity: HIGH\ncount: 3\ntags[2]: x,y\n"
	if got := buf.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestToonPrinter_EmptySlice(t *testing.T) {
	var buf bytes.Buffer
	if err := NewPrinter(FormatTOON, nil).Print(&buf, []toonItem{}); err != nil {
		t.Fatalf("Print: %v", err)
	}
	if got := buf.String(); got != "[]\n" {
		t.Errorf("got %q", got)
	}
}

func TestToonPrinter_LargeIntegerIsNotMangled(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]int64{"n": 1234567890123}
	if err := NewPrinter(FormatTOON, nil).Print(&buf, data); err != nil {
		t.Fatalf("Print: %v", err)
	}
	if got := buf.String(); got != "n: 1234567890123\n" {
		t.Errorf("got %q", got)
	}
}

type toonRisk struct {
	ID        string     `json:"id"`
	Severity  string     `json:"severity"`
	Status    string     `json:"status"`
	Provider  string     `json:"provider"`
	Region    string     `json:"region"`
	Resource  string     `json:"resource_type"`
	Score     int64      `json:"score"`
	Tags      []string   `json:"tags"`
	Owner     toonOwner  `json:"owner"`
	FirstSeen string     `json:"first_seen"`
	Notes     []toonNote `json:"notes,omitempty"`
}

type toonOwner struct {
	Team  string `json:"team"`
	Email string `json:"email"`
}

type toonNote struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

func toonRiskFixture(n int) []toonRisk {
	severities := []string{"CRITICAL", "HIGH", "MEDIUM", "LOW"}
	providers := []string{"aws", "azure", "gcp"}
	resources := []string{"AWS::S3::Bucket", "AWS::EC2::Instance", "Microsoft.Compute/virtualMachines", "google_compute_instance"}
	risks := make([]toonRisk, n)
	for i := range risks {
		risks[i] = toonRisk{
			ID:        fmt.Sprintf("risk-%04d-abcdef", i),
			Severity:  severities[i%len(severities)],
			Status:    "open",
			Provider:  providers[i%len(providers)],
			Region:    "us-east-1",
			Resource:  resources[i%len(resources)],
			Score:     int64(40 + i%60),
			Tags:      []string{"prod", "internet-facing"},
			Owner:     toonOwner{Team: "platform", Email: "platform@example.com"},
			FirstSeen: "2026-09-01T12:00:00Z",
		}
	}
	return risks
}

// approxTokens is a dependency-free estimate of LLM token count: each
// alphanumeric run costs ceil(len/4) tokens, each punctuation/symbol rune costs
// one, and whitespace is free. It is only meant for comparing formats.
func approxTokens(s string) int {
	tokens, run := 0, 0
	flush := func() {
		tokens += (run + 3) / 4
		run = 0
	}
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			run++
		case unicode.IsSpace(r):
			flush()
		default:
			flush()
			tokens++
		}
	}
	flush()
	return tokens
}

// TestToonTokenSavings prints a size/token comparison of --output json vs
// --output toon. Run with -v to see the report.
func TestToonTokenSavings(t *testing.T) {
	data := toonRiskFixture(25)

	render := func(f Format) string {
		var buf bytes.Buffer
		if err := NewPrinter(f, nil).Print(&buf, data); err != nil {
			t.Fatalf("Print(%s): %v", f, err)
		}
		return buf.String()
	}

	jsonOut, toonOut := render(FormatJSON), render(FormatTOON)
	jsonTok, toonTok := approxTokens(jsonOut), approxTokens(toonOut)

	t.Logf("fixture: %d risks", len(data))
	t.Logf("%-6s %7d bytes  ~%6d tokens", "json", len(jsonOut), jsonTok)
	t.Logf("%-6s %7d bytes  ~%6d tokens", "toon", len(toonOut), toonTok)
	t.Logf("savings: %.1f%% bytes, ~%.1f%% tokens (estimate, not a real tokenizer)",
		100*float64(len(jsonOut)-len(toonOut))/float64(len(jsonOut)),
		100*float64(jsonTok-toonTok)/float64(jsonTok))

	if toonTok >= jsonTok {
		t.Errorf("expected TOON (%d tokens) to use fewer tokens than JSON (%d)", toonTok, jsonTok)
	}
}
