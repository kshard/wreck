package wreck_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/kshard/wreck"
)

func TestStreamJSON(t *testing.T) {
	var b bytes.Buffer
	input := []byte("Hello World!\n")

	w := wreck.NewWriterJSON(&b)
	if _, err := w.Write(input); err != nil {
		t.Errorf("unable to write JSON stream: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Errorf("unable to close JSON stream: %v", err)
	}

	r, err := wreck.NewReaderJSON(&b)
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
