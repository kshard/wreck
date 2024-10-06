//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.
// https://github.com/kshard/wreck
//

package wreck_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/kshard/wreck"
)

func TestStreamJSONWithGzip(t *testing.T) {
	var b bytes.Buffer
	input := []byte("Hello World!\n")

	w := wreck.NewWriterJSON(&b, true)
	if _, err := w.Write(input); err != nil {
		t.Errorf("unable to write JSON stream: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Errorf("unable to close JSON stream: %v", err)
	}

	r, err := wreck.NewReaderJSON(&b, true)
	if err != nil {
		t.Errorf("unable to create JSON stream: %v", err)
	}

	d, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("unable to read JSON stream: %v", err)
	}

	if !bytes.Equal(d, input) {
		t.Errorf("stream corrupted: got %v, expected %v", d, input)
	}
}

func TestStreamJSONWithoutGzip(t *testing.T) {
	var b bytes.Buffer
	input := []byte("Hello World!\n")

	w := wreck.NewWriterJSON(&b, false)
	if _, err := w.Write(input); err != nil {
		t.Errorf("unable to write JSON stream: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Errorf("unable to close JSON stream: %v", err)
	}

	r, err := wreck.NewReaderJSON(&b, false)
	if err != nil {
		t.Errorf("unable to create JSON stream: %v", err)
	}

	d, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("unable to read JSON stream: %v", err)
	}

	if !bytes.Equal(d, input) {
		t.Errorf("stream corrupted: got %v, expected %v", d, input)
	}
}
