package wreck_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/kshard/wreck"
)

func TestCodec(t *testing.T) {
	sizes := []int{0, 1, 3, 16, 256, 1024}

	for _, szData := range sizes {
		for _, szUniqueKey := range sizes {
			for _, szSortKey := range sizes {
				input := wreck.Chunk{
					UniqueKey: chunk(szUniqueKey),
					SortKey:   chunk(szSortKey),
					Vector:    chunk(szData),
				}

				b := &bytes.Buffer{}
				if err := wreck.Encode(b, &input); err != nil {
					t.Errorf("unable to encode chunk: %v", err)
				}

				var output wreck.Chunk
				if err := wreck.Decode(b, &output); err != nil {
					t.Errorf("unable to decode chunk: %v", err)
				}

				if !bytes.Equal(input.UniqueKey, output.UniqueKey) {
					t.Errorf("unique key is corrupted")
				}

				if !bytes.Equal(input.SortKey, output.SortKey) {
					t.Errorf("sort key is corrupted")
				}

				if !bytes.Equal(input.Vector, output.Vector) {
					t.Errorf("data key is corrupted")
				}
			}
		}
	}

}

func chunk(n int) []uint8 {
	if n == 0 {
		return nil
	}

	c := make([]uint8, n)
	for i := 0; i < n; i++ {
		c[i] = uint8(rand.Int31n(256))
	}
	return c
}
