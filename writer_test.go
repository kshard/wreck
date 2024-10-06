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

func TestWriter(t *testing.T) {
	sizes := []int{0, 1, 3, 16, 256, 1024}

	for _, szData := range sizes {
		for _, szUniqueKey := range sizes {
			for _, szSortKey := range sizes {
				input := data(szUniqueKey, szSortKey, szData)

				b := &bytes.Buffer{}
				w := wreck.NewWriter[uint8](b)

				if err := w.Write(input.UniqueKey, input.SortKey, input.Vector); err != nil {
					t.Errorf("unable to encode chunk: %v", err)
				}

				assertEncodedChunk(t, b, input)
			}
		}
	}
}

func TestEncoder(t *testing.T) {
	sizes := []int{0, 1, 3, 16, 256, 1024}

	w := wreck.NewEncoder[uint8]()

	for _, szData := range sizes {
		for _, szUniqueKey := range sizes {
			for _, szSortKey := range sizes {
				input := data(szUniqueKey, szSortKey, szData)

				pckt, err := w.Encode(input.UniqueKey, input.SortKey, input.Vector)
				if err != nil {
					t.Errorf("unable to encode chunk: %v", err)
				}

				assertEncodedChunk(t, bytes.NewBuffer(pckt), input)
			}
		}
	}
}

func assertEncodedChunk(t *testing.T, r io.Reader, expected pack) {
	t.Helper()

	var output wreck.Chunk
	if err := wreck.Decode(r, &output); err != nil {
		t.Errorf("unable to decode chunk: %v", err)
	}

	if !bytes.Equal(expected.UniqueKey, output.UniqueKey) {
		t.Errorf("unique key is corrupted")
	}

	if !bytes.Equal(expected.SortKey, output.SortKey) {
		t.Errorf("sort key is corrupted")
	}

	if !bytes.Equal(expected.Vector, output.Vector) {
		t.Errorf("data key is corrupted")
	}
}
