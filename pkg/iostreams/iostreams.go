// Copyright (c) 2022 CrowdStrike, Inc.
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

package iostreams

import (
	"io"
	"os"

	"golang.org/x/term"
)

type IOStreams struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer

	stdinIsTTY  bool
	stdoutIsTTY bool
	stderrIsTTY bool

	neverPrompt bool
}

func (s *IOStreams) NewIOStreams() *IOStreams {
	io := &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	stdoutIsTTy := isTerminal(os.Stdout)
	stderrIsTTy := isTerminal(os.Stderr)

	io.SetStdoutTTY(stdoutIsTTy)
	io.SetStderrTTY(stderrIsTTy)

	return io
}

func (s *IOStreams) SetStdinTTY(tty bool) {
	s.stdinIsTTY = tty
}

func (s *IOStreams) SetStdoutTTY(tty bool) {
	s.stdoutIsTTY = tty
}

func (s *IOStreams) SetStderrTTY(tty bool) {
	s.stderrIsTTY = tty
}

func (s *IOStreams) StdinTTY() bool {
	return s.stdinIsTTY
}

func (s *IOStreams) CanPrompt() bool {
	if s.neverPrompt {
		return false
	}
	return s.IsStdinTTY() && s.IsStdoutTTY()
}

func (s *IOStreams) GetNeverPrompt() bool {
	return s.neverPrompt
}

func (s *IOStreams) SetNeverPrompt(neverPrompt bool) {
	s.neverPrompt = neverPrompt
}

func (s *IOStreams) IsStderrTTY() bool {
	if stderr, ok := s.ErrOut.(*os.File); ok {
		return isTerminal(stderr)
	}
	return false
}

func (s *IOStreams) IsStdoutTTY() bool {
	if stdout, ok := s.Out.(*os.File); ok {
		return isTerminal(stdout)
	}
	return false
}

func (s *IOStreams) IsStdinTTY() bool {
	if stdin, ok := s.In.(*os.File); ok {
		return isTerminal(stdin)
	}
	return false
}

func isTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}
