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
	"math"
)

// Stream scanner
type Scanner[T any] struct {
	r    io.Reader
	b    Chunk
	err  error
	fmap func([]uint8) []T
}

func NewScanner[T any](r io.Reader) *Scanner[T] {
	codec := &Scanner[T]{r: r}

	switch any(*new(T)).(type) {
	case float32:
		codec.fmap = toFloat32
	case uint8:
		codec.fmap = toUint8
	}

	return codec
}

func (codec *Scanner[T]) WithDecoder(fmap func([]uint8) []T) {
	codec.fmap = fmap
}

func (codec *Scanner[T]) Err() error         { return codec.err }
func (codec *Scanner[T]) UniqueKey() []uint8 { return codec.b.UniqueKey }
func (codec *Scanner[T]) SortKey() []uint8   { return codec.b.SortKey }
func (codec *Scanner[T]) Vector() []T        { return codec.fmap(codec.b.Vector) }

func (codec *Scanner[T]) Scan() bool {
	if err := Decode(codec.r, &codec.b); err != nil {
		if !errors.Is(err, io.EOF) {
			codec.err = err
			return false
		}
		return false
	}

	return true
}

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
