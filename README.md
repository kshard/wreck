# wreck

Wreck is a binary format for efficient streaming, storing and exchanging "vectors" that maximizes efficiency and minimizes overhead. The format annotates "vectors" with unique and sort keys making it possible to store and lookup without doing a full deserialization. 


## Principles

**Compactness**: Aim to minimize the size of the serialized data. This involves using efficient encoding schemes, and minimizing metadata overhead.

**Simplicity**: The format is easy to parse and generate using streaming codecs. It avoids overly complex structures and uses only byte sequences.

**Efficiency**: Ensure that both serialization and deserialization processes are fast. Use fixed-size data types where possible to speed up parsing, and minimize the need for complex computations or lookups.

**Cross-platform Compatibility**: The format uses only little endian encoding and octet-streams. It avoids any character encodings. Floats are encoded using IEEE 754 binary representation. It ensures the correctness across different platforms and architectures. 

The format simplify the implementation through establishing dependencies to external stream codecs in the following aspects. It makes a limitation that codec cannot be used standalone and requires applications to negotiate these parameters.  

**Security**: Use external stream ciphers.

**Integrity**: Use external streaming error detection and error correction schemas.

**Compression**: Use external compression. The nature of the data does not allowing extreme gains with compression. Gzip saves only 16% with best compression.

## Format

```
//
// 0x00 : 4 byte     | Block Size    (L)
// 0x04 : 2 byte     | Vector Size   (V)
// 0x06 : 2 byte     | Sort Key Size (S)
// 0x08 : 1 byte * V | Vector
// 0xZZ : 1 byte * S | Sort Key
// 0xXX : 1 byte *   | Unique Key
//
```


## Getting started

The latest version of the module is available at `main` branch. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

Use `go get` to retrieve the library and add it as dependency to your application.

```bash
go get -u github.com/kshard/wreck
```

### Quick Example

```go
// Create writer for []float32 vector
w := wreck.NewWriter[float32](out)

// Writer vector
if err := w.Write(uniqueKey, sortKey, vector); err != nil {
  // ...
}
```

```go
// Create scanner for []float32 vector
r := wreck.NewScanner[float32](in)

// Scan through stream
for r.Scan() {
  // consume vector  
  r.UniqueKey()
  r.Vector()
}

if err := r.Err(); err != nil {
  // ...
}
```

### Use-cases

* Large vector streams
  - Output `Writer[T any]`
  - Input `Scanner[T any]`
* Batching vectors, using JSON as primary protocol
  - On-the-wire protocol encoding/decoding with `WriterJSON` and `ReaderJSON`
  - Output `Writer[T any]`
  - Input `Scanner[T any]` 
* Transmitting one vector in the packet
  - Output `Encoder[T any]`
  - Input `Decoder[T any]`


## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

The build and testing process requires [Go](https://golang.org) version 1.21 or later.


### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/kshard/wreck/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/kshard/wreck.svg?style=for-the-badge)](LICENSE)

