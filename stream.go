//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.
// https://github.com/kshard/wreck
//

package wreck

import (
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

// Output stream for of vectors to be used as [json.RawMessage].
// The stream is base64 encoded string "H4sI...eQCAN" compatible with
// Golang's []byte JSON encoding protocol.
type WriterJSON []io.WriteCloser

// Create output stream for of vectors to be used as [json.RawMessage].
func NewWriterJSON(w io.Writer, withGzip bool) io.WriteCloser {
	q := newQuoter(w)
	b := base64.NewEncoder(base64.StdEncoding, q)
	if !withGzip {
		return WriterJSON{b, q}
	}

	z := gzip.NewWriter(b)

	return WriterJSON{z, b, q}
}

func (w WriterJSON) Close() error {
	for _, c := range w {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (w WriterJSON) Write(p []byte) (n int, err error) {
	return w[0].Write(p)
}

// Create input stream to read vectors from [json.RawMessage].
func NewReaderJSON(r io.Reader, withGzip bool) (io.Reader, error) {
	q, err := newUnQuote(r)
	if err != nil {
		return nil, err
	}

	b := base64.NewDecoder(base64.StdEncoding, q)
	if !withGzip {
		return b, nil
	}

	z, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}

	return z, nil
}

// Quoter stream escapes stream content into double quotes (")
type quoter struct {
	w     io.Writer
	start bool
}

func newQuoter(w io.Writer) io.WriteCloser {
	return &quoter{w: w}
}

func (w *quoter) Close() error {
	if _, err := w.w.Write([]byte{'"'}); err != nil {
		return err
	}

	return nil
}

func (w *quoter) Write(p []byte) (n int, err error) {
	if !w.start {
		if _, err := w.w.Write([]byte{'"'}); err != nil {
			return 0, err
		}
		w.start = true
	}

	return w.w.Write(p)
}

// UnQuoter stream unescape stream content into double quotes (")
type unquote struct {
	r   io.Reader
	end bool
}

func newUnQuote(r io.Reader) (io.Reader, error) {
	b := []byte{0x00}
	if _, err := r.Read(b); err != nil {
		return nil, err
	}

	if b[0] != '"' {
		return nil, fmt.Errorf("invalid format")
	}

	return &unquote{r: r}, nil
}

func (r *unquote) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	if err != nil && err != io.EOF {
		return n, err
	}

	// Skip the last `"`
	if !r.end {
		if n > 0 && p[n-1] == '"' {
			n--
		}
		r.end = true
	}

	return n, err
}
