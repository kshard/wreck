//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.
// https://github.com/kshard/wreck
//

package wreck

import (
	"encoding/binary"
	"errors"
	"io"
)

type Chunk struct {
	UniqueKey []uint8
	SortKey   []uint8
	Vector    []uint8
}

func Encode(w io.Writer, wreck *Chunk) error {
	var head [8]uint8

	szData := len(wreck.Vector)
	szUniqueKey := len(wreck.UniqueKey)
	szSortKey := len(wreck.SortKey)
	szWreck := szData + szUniqueKey + szSortKey

	binary.LittleEndian.PutUint32(head[0:4], uint32(szWreck))
	binary.LittleEndian.PutUint16(head[4:6], uint16(szData))
	binary.LittleEndian.PutUint16(head[6:8], uint16(szSortKey))

	if _, err := w.Write(head[:]); err != nil {
		return err
	}

	if szData > 0 {
		if _, err := w.Write(wreck.Vector); err != nil {
			return err
		}
	}

	if szSortKey > 0 {
		if _, err := w.Write(wreck.SortKey); err != nil {
			return err
		}
	}

	if szUniqueKey > 0 {
		if _, err := w.Write(wreck.UniqueKey); err != nil {
			return err
		}
	}

	return nil
}

func Decode(r io.Reader, wreck *Chunk) error {
	var head [8]uint8

	if _, err := io.ReadFull(r, head[:]); err != nil {
		return err
	}

	szWreck := int(binary.LittleEndian.Uint32(head[0:4]))
	szData := int(binary.LittleEndian.Uint16(head[4:6]))
	szSortKey := int(binary.LittleEndian.Uint16(head[6:8]))
	szUniqueKey := szWreck - szData - szSortKey

	b := make([]uint8, szWreck)
	if _, err := io.ReadFull(r, b); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
		return io.ErrUnexpectedEOF
	}

	if szData == 0 {
		wreck.Vector = nil
	} else {
		wreck.Vector = b[0:szData]
	}

	if szSortKey == 0 {
		wreck.SortKey = nil
	} else {
		wreck.SortKey = b[szData : szData+szSortKey]
	}

	if szUniqueKey == 0 {
		wreck.UniqueKey = nil
	} else {
		wreck.UniqueKey = b[szData+szSortKey : szData+szSortKey+szUniqueKey]
	}

	return nil
}
