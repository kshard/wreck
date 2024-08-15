package wreck_test

import (
	"bytes"
	"testing"

	"github.com/kshard/wreck"
)

func TestWriter(t *testing.T) {
	sizes := []int{0, 1, 3, 16, 256, 1024}

	for _, szData := range sizes {
		for _, szUniqueKey := range sizes {
			for _, szSortKey := range sizes {
				uniqueKey := chunk(szUniqueKey)
				sortKey := chunk(szSortKey)
				vector := chunk(szData)

				b := &bytes.Buffer{}
				w := wreck.NewWriter[uint8](b)

				if err := w.Write(uniqueKey, sortKey, vector); err != nil {
					t.Errorf("unable to encode chunk: %v", err)
				}

				var output wreck.Chunk
				if err := wreck.Decode(b, &output); err != nil {
					t.Errorf("unable to decode chunk: %v", err)
				}

				if !bytes.Equal(uniqueKey, output.UniqueKey) {
					t.Errorf("unique key is corrupted")
				}

				if !bytes.Equal(sortKey, output.SortKey) {
					t.Errorf("sort key is corrupted")
				}

				if !bytes.Equal(vector, output.Vector) {
					t.Errorf("data key is corrupted")
				}
			}
		}
	}
}
