# tar
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/tar)](https://pkg.go.dev/github.com/hslam/tar)
[![Build Status](https://github.com/hslam/tar/workflows/build/badge.svg)](https://github.com/hslam/tar/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/tar)](https://goreportcard.com/report/github.com/hslam/tar)
[![LICENSE](https://img.shields.io/github/license/hslam/tar.svg?style=flat-square)](https://github.com/hslam/tar/blob/master/LICENSE)

Package tar implements access to tar archives.

## Feature
* Tar
* Gzip
* Tar files or dirs

## Get started

### Install
```
go get github.com/hslam/tar
```
### Import
```
import "github.com/hslam/tar"
```
### Usage
#### Example
```go
package main

import (
	"fmt"
	"github.com/hslam/tar"
	"os"
)

func main() {
	name := "file"
	targz := "file.tar.gz"
	defer os.Remove(name)
	defer os.Remove(targz)
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	contents := "Hello World"
	file.Write([]byte(contents))
	file.Close()
	tar.Targz(targz, name)
	os.Remove(name)
	tar.Untargz(targz)
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := make([]byte, len(contents))
	f.Read(buf)
	fmt.Println(string(buf))
}
```

#### Tar bytes example
```go
package main

import (
	"fmt"
	"github.com/hslam/tar"
	"os"
)

func main() {
	targz := "file.tar.gz"
	defer os.Remove(targz)
	tw, err := tar.NewGzipFileWriter(targz)
	if err != nil {
		panic(err)
	}
	tw.TarBytes("file", []byte("Hello World"))
	tw.Flush()
	tw.Close()
	tr, err := tar.NewGzipFileReader(targz)
	if err != nil {
		panic(err)
	}
	_, data, err := tr.NextBytes()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
```

#### Output
```
Hello World
```

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)

### Author
tar was written by Meng Huang.


