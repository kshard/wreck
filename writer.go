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
	"io"
	"math"
)

// Stream writer
type Writer[T any] struct {
	w    io.Writer
	b    Chunk
	fmap func([]T) []uint8
}

func NewWriter[T any](w io.Writer) *Writer[T] {
	codec := &Writer[T]{w: w}

	switch any(*new(T)).(type) {
	case float32:
		codec.fmap = fromFloat32
	case uint8:
		codec.fmap = fromInt8
	}

	return codec
}

func (codec *Writer[T]) WithEncoder(fmap func([]T) []uint8) {
	codec.fmap = fmap
}

func (codec *Writer[T]) Write(uniqueKey []uint8, sortKey []uint8, vector []T) error {
	codec.b.UniqueKey = uniqueKey
	codec.b.SortKey = sortKey
	codec.b.Vector = codec.fmap(vector)

	return Encode(codec.w, &codec.b)
}

func fromFloat32[T any](v []T) []uint8 {
	vv := any(v).([]float32)

	p := 0
	b := make([]uint8, len(vv)*4)
	for i := 0; i < len(b); i += 4 {
		f := math.Float32bits(vv[p])
		binary.LittleEndian.PutUint32(b[i:i+4], f)
		p++
	}

	return b
}

func fromInt8[T any](v []T) []uint8 {
	return any(v).([]uint8)
}
