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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/toon-format/toon-go"
)

// toonPrinter writes data as TOON (Token-Oriented Object Notation), a compact
// encoding of the JSON data model that uses fewer tokens when fed to LLMs.
//
// Values are first marshalled with encoding/json so the content is identical to
// --output json (json tags, omitempty, custom marshallers), then re-encoded as
// TOON with field order preserved.
type toonPrinter struct{}

func (p *toonPrinter) Print(w io.Writer, data any) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshalling to JSON: %w", err)
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	v, err := decodeOrdered(dec)
	if err != nil {
		return fmt.Errorf("converting to TOON: %w", err)
	}

	s, err := toon.MarshalString(v)
	if err != nil {
		return fmt.Errorf("encoding TOON: %w", err)
	}
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	_, err = io.WriteString(w, s)
	return err
}

// decodeOrdered reads one JSON value from dec, turning objects into
// toon.Object (to keep key order) and numbers into int64 or float64.
func decodeOrdered(dec *json.Decoder) (any, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}

	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			var fields []toon.Field
			for dec.More() {
				keyTok, err := dec.Token()
				if err != nil {
					return nil, err
				}
				key, ok := keyTok.(string)
				if !ok {
					return nil, errors.New("non-string object key")
				}
				val, err := decodeOrdered(dec)
				if err != nil {
					return nil, err
				}
				fields = append(fields, toon.Field{Key: key, Value: val})
			}
			if _, err := dec.Token(); err != nil { // closing '}'
				return nil, err
			}
			return toon.NewObject(fields...), nil
		case '[':
			arr := []any{}
			for dec.More() {
				val, err := decodeOrdered(dec)
				if err != nil {
					return nil, err
				}
				arr = append(arr, val)
			}
			if _, err := dec.Token(); err != nil { // closing ']'
				return nil, err
			}
			return arr, nil
		}
		return nil, fmt.Errorf("unexpected delimiter %q", t)
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return i, nil
		}
		f, err := t.Float64()
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return tok, nil
	}
}
