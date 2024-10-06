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
	"testing"

	"github.com/kshard/wreck"
)

func TestScanner(t *testing.T) {
	sizes := []int{0, 1, 3, 16, 256, 1024}

	for _, szData := range sizes {
		for _, szUniqueKey := range sizes {
			for _, szSortKey := range sizes {
				input := data(szUniqueKey, szSortKey, szData)

				b := &bytes.Buffer{}
				w := wreck.NewWriter[uint8](b)
				r := wreck.NewScanner[uint8](b)

				if err := w.Write(input.UniqueKey, input.SortKey, input.Vector); err != nil {
					t.Errorf("unable to encode chunk: %v", err)
				}

				if has := r.Scan(); !has {
					t.Errorf("unable to decode chunk: unexpected eof")
				}
				if err := r.Err(); err != nil {
					t.Errorf("unable to decode chunk: %v", err)
				}

				if !bytes.Equal(input.UniqueKey, r.UniqueKey()) {
					t.Errorf("unique key is corrupted")
				}

				if !bytes.Equal(input.SortKey, r.SortKey()) {
					t.Errorf("sort key is corrupted")
				}

				if !bytes.Equal(input.Vector, r.Vector()) {
					t.Errorf("data key is corrupted")
				}
			}
		}
	}
}
