# bencode

[![Build Status][build-img]][build-url]
[![GoDoc][pkg-img]][pkg-url]
[![Go Report Card][reportcard-img]][reportcard-url]
[![Coverage][coverage-img]][coverage-url]

Package implements Bencode encoding and decoding in Go.

## Features

* Simple API.
* Clean and tested code.
* Optimized for speed.
* Dependency-free.

## Install

Go version 1.15+

```
go get github.com/cristalhq/bencode
```

## Example

```go
var data yourStruct{}
buf, err := bencode.Marshal(&data)
checkErr(err)

// or via Encoder:
w := &bytes.Buffer{} // or any other io.Writer
err = bencode.NewEncoder(w).Encode(&data)
checkErr(err)

err = bencode.Unmarshal(buf, &data)
checkErr(err)

// or via Decoder:
r := &bytes.NewBufferString("...") // or any other io.Reader
err = bencode.NewDecoder(r).Decode(&data)
checkErr(err)
```

## Documentation

See [these docs][pkg-url].

## License

[MIT License](LICENSE).

[build-img]: https://github.com/cristalhq/bencode/workflows/build/badge.svg
[build-url]: https://github.com/cristalhq/bencode/actions
[pkg-img]: https://pkg.go.dev/badge/cristalhq/bencode
[pkg-url]: https://pkg.go.dev/github.com/cristalhq/bencode
[reportcard-img]: https://goreportcard.com/badge/cristalhq/bencode
[reportcard-url]: https://goreportcard.com/report/cristalhq/bencode
[coverage-img]: https://codecov.io/gh/cristalhq/bencode/branch/master/graph/badge.svg
[coverage-url]: https://codecov.io/gh/cristalhq/bencode
