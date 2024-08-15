package main

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/kshard/wreck"
)

func main() {
	b := &bytes.Buffer{}

	//
	// Writing
	//

	// create compressed wreck writer stream
	e := gzip.NewWriter(b)
	w := wreck.NewWriter[float32](e)

	// write vectors
	for i := uint8(0); i < 255; i++ {
		// write only UniqueKey and Vector
		err := w.Write([]uint8{i}, nil, []float32{0.1, 0.2, 0.3, 0.4})
		if err != nil {
			panic(err)
		}
	}
	e.Close()

	//
	// Reading
	//

	// create compressed wreck reader stream
	d, err := gzip.NewReader(b)
	if err != nil {
		panic(err)
	}
	r := wreck.NewScanner[float32](d)

	for r.Scan() {
		fmt.Printf("%x : %v\n", r.UniqueKey(), r.Vector())
	}

	if err := r.Err(); err != nil {
		panic(err)
	}

	d.Close()
}
