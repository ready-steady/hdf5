# HDF5 [![Build Status][travis-img]][travis-url]

The package provides a reader and writer of [HDF5][1] files.

## [Documentation][doc]

## Installation

Fetch the package:

```bash
go get -d github.com/ready-steady/hdf5
```

Go to the directory of the package:

```bash
cd $GOPATH/src/github.com/ready-steady/hdf5
```

Install the package:

```bash
make install
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/ready-steady/hdf5"
)

func main() {
	put("data.h5")
	get("data.h5")
}

func put(path string) {
	file, _ := hdf5.Create(path)
	defer file.Close()

	A := 42
	file.Put("A", A)

	B := []float64{1, 2, 3}
	file.Put("B", B)

	C := struct {
		D int
		E []float64
	}{
		D: 42,
		E: []float64{1, 2, 3},
	}
	file.Put("C", C)
}

func get(path string) {
	file, _ := hdf5.Open(path)
	defer file.Close()

	A := 0
	file.Get("A", &A)
	fmt.Println(A)

	B := []float64{}
	file.Get("B", &B)
	fmt.Println(B)

	C := struct {
		D int
		E []float64
	}{}
	file.Get("C", &C)
	fmt.Println(C)
}
```

## Contribution

1. Fork the project.
2. Implement your idea.
3. Open a pull request.

[1]: https://en.wikipedia.org/wiki/Hierarchical_Data_Format

[doc]: http://godoc.org/github.com/ready-steady/hdf5
[travis-img]: https://travis-ci.org/ready-steady/hdf5.svg?branch=master
[travis-url]: https://travis-ci.org/ready-steady/hdf5
