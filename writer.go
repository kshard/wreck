//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.
// https://github.com/kshard/wreck
//

package wreck

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// basic vector encoder
type encoder[T any] struct {
	chunk        Chunk
	fmap         func([]T) []uint8
	maxUniqueKey int
	maxSortKey   int
	maxVector    int
}

func newEncoder[T any]() encoder[T] {
	codec := encoder[T]{}

	switch any(*new(T)).(type) {
	case float32:
		codec.fmap = fromFloat32
	case uint8:
		codec.fmap = fromInt8
	}

	return codec
}

func (codec *encoder[T]) WithEncoder(fmap func([]T) []uint8) {
	codec.fmap = fmap
}

func (codec *encoder[T]) WithMaxUniqueKey(v int) {
	codec.maxUniqueKey = v
}

func (codec *encoder[T]) WithMaxSortKey(v int) {
	codec.maxSortKey = v
}

func (codec *encoder[T]) WithMaxVector(v int) {
	codec.maxVector = v
}

func (codec *encoder[T]) encode(uniqueKey, sortKey []uint8, vec []T) error {
	if codec.maxUniqueKey > 0 && len(uniqueKey) > codec.maxUniqueKey {
		return fmt.Errorf("length exceeded : uniqueKey (%d)", len(uniqueKey))
	}
	if codec.maxSortKey > 0 && len(sortKey) > codec.maxSortKey {
		return fmt.Errorf("length exceeded : sortKey (%d)", len(sortKey))
	}
	if codec.maxVector > 0 && len(vec) > codec.maxVector {
		return fmt.Errorf("length exceeded : vector (%d)", len(vec))
	}

	codec.chunk.UniqueKey = uniqueKey
	codec.chunk.SortKey = sortKey
	codec.chunk.Vector = codec.fmap(vec)

	return nil
}

//------------------------------------------------------------------------------

// Binary stream vector writer
type Writer[T any] struct {
	encoder[T]
	w io.Writer
}

// Create instance of vector writer
func NewWriter[T any](w io.Writer) *Writer[T] {
	return &Writer[T]{
		encoder: newEncoder[T](),
		w:       w,
	}
}

// Write vectors
func (codec *Writer[T]) Write(uniqueKey, sortKey []uint8, vec []T) error {
	if err := codec.encode(uniqueKey, sortKey, vec); err != nil {
		return err
	}

	return Encode(codec.w, &codec.chunk)
}

//------------------------------------------------------------------------------

// Encode single vector as binary packet
type Encoder[T any] struct{ encoder[T] }

func NewEncoder[T any]() *Encoder[T] {
	return &Encoder[T]{
		encoder: newEncoder[T](),
	}
}

func (codec *Encoder[T]) Encode(uniqueKey, sortKey []uint8, vec []T) ([]byte, error) {
	if err := codec.encode(uniqueKey, sortKey, vec); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := Encode(&buf, &codec.chunk); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//------------------------------------------------------------------------------

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
