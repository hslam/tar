# tar
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
	tname := name + tar.TargzSuffix
	defer os.Remove(name)
	defer os.Remove(tname)
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	contexts := "Hello World"
	file.Write([]byte(contexts))
	file.Close()
	tar.Targz(tname, name)
	os.Remove(name)
	tar.Untargz(tname)
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := make([]byte, len(contexts))
	f.Read(buf)
	fmt.Println(string(buf))
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


