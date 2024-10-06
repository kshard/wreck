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
	"errors"
	"fmt"
	"io"
	"math"
)

// basic vector decoder
type decoder[T any] struct {
	fmap         func([]uint8) []T
	maxUniqueKey int
	maxSortKey   int
	maxVector    int
}

func newDecoder[T any]() decoder[T] {
	codec := decoder[T]{}

	switch any(*new(T)).(type) {
	case float32:
		codec.fmap = toFloat32
	case uint8:
		codec.fmap = toUint8
	}

	return codec
}

func (codec *decoder[T]) WithDecoder(fmap func([]uint8) []T) {
	codec.fmap = fmap
}

func (codec *decoder[T]) WithMaxUniqueKey(v int) {
	codec.maxUniqueKey = v
}

func (codec *decoder[T]) WithMaxSortKey(v int) {
	codec.maxSortKey = v
}

func (codec *decoder[T]) WithMaxVector(v int) {
	codec.maxVector = v
}

func (codec *decoder[T]) decode(chunk *Chunk) ([]uint8, []uint8, []T, error) {
	uniqueKey := chunk.UniqueKey
	sortKey := chunk.SortKey
	vec := codec.fmap(chunk.Vector)

	if codec.maxUniqueKey > 0 && len(uniqueKey) > codec.maxUniqueKey {
		return nil, nil, nil, fmt.Errorf("length exceeded : uniqueKey (%d)", len(uniqueKey))
	}
	if codec.maxSortKey > 0 && len(sortKey) > codec.maxSortKey {
		return nil, nil, nil, fmt.Errorf("length exceeded : sortKey (%d)", len(sortKey))
	}
	if codec.maxVector > 0 && len(vec) > codec.maxVector {
		return nil, nil, nil, fmt.Errorf("length exceeded : vector (%d)", len(vec))
	}

	return uniqueKey, sortKey, vec, nil
}

//------------------------------------------------------------------------------

// Vector stream
type Scanner[T any] struct {
	decoder[T]
	r         io.Reader
	uniqueKey []uint8
	sortKey   []uint8
	vec       []T

	err error
}

func NewScanner[T any](r io.Reader) *Scanner[T] {
	return &Scanner[T]{
		decoder: newDecoder[T](),
		r:       r,
	}
}

func (codec *Scanner[T]) Err() error         { return codec.err }
func (codec *Scanner[T]) UniqueKey() []uint8 { return codec.uniqueKey }
func (codec *Scanner[T]) SortKey() []uint8   { return codec.sortKey }
func (codec *Scanner[T]) Vector() []T        { return codec.vec }

func (codec *Scanner[T]) Scan() bool {
	var chunk Chunk

	if err := Decode(codec.r, &chunk); err != nil {
		if !errors.Is(err, io.EOF) {
			codec.err = err
			return false
		}
		return false
	}

	codec.uniqueKey, codec.sortKey, codec.vec, codec.err = codec.decode(&chunk)
	return codec.err == nil
}

//------------------------------------------------------------------------------

// Decode binary packet to vector
type Decoder[T any] struct{ decoder[T] }

func NewDecoder[T any]() *Decoder[T] {
	return &Decoder[T]{
		decoder: newDecoder[T](),
	}
}

func (codec *Decoder[T]) Decode(pack []byte, uniqueKey, sortKey *[]uint8, vec *[]T) (err error) {
	var chunk Chunk

	if err := Decode(bytes.NewBuffer(pack), &chunk); err != nil {
		return err
	}

	*uniqueKey, *sortKey, *vec, err = codec.decode(&chunk)
	return nil
}

//------------------------------------------------------------------------------

func toFloat32[T any](b []uint8) []T {
	v := make([]float32, len(b)/4)

	p := 0
	for i := 0; i < len(b); i += 4 {
		v[p] = math.Float32frombits(binary.LittleEndian.Uint32(b[i : i+4]))
		p++
	}

	return any(v).([]T)
}

func toUint8[T any](b []uint8) []T {
	return any(b).([]T)
}
