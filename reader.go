// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package tar

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ErrTooManyArgs is returned when the arguments length is bigger than 1.
var ErrTooManyArgs = errors.New("too many arguments")

// Reader provides sequential access to the contents of a tar archive.
type Reader struct {
	*tar.Reader
	gr *gzip.Reader
	f  *os.File
}

// NewReader creates a new Reader reading from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{Reader: tar.NewReader(r)}
}

// NewGzipReader creates a new gzip Reader reading from r.
func NewGzipReader(r io.Reader) (*Reader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &Reader{Reader: tar.NewReader(gr), gr: gr}, nil
}

// NewFileReader creates a new Reader reading from file.
func NewFileReader(name string) (*Reader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &Reader{Reader: tar.NewReader(f), f: f}, nil
}

// NewGzipFileReader creates a new gzip Reader reading from file.
func NewGzipFileReader(name string) (r *Reader, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(f)
	if err == nil {
		r = &Reader{Reader: tar.NewReader(gr), gr: gr, f: f}
	}
	return r, err
}

// Untar untars all the files to dir.
func (t *Reader) Untar(dir ...string) (files, dirs []string, err error) {
	var path string
	var isDir bool
	for {
		path, isDir, err = t.NextFile(dir...)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return files, dirs, err
		}
		if len(path) > 0 {
			if isDir {
				dirs = append(dirs, path)
			} else {
				files = append(files, path)
			}
		}
	}
}

// NextFile advances to the next file in the tar archive.
func (t *Reader) NextFile(dir ...string) (name string, isDir bool, err error) {
	var dirpath string
	if len(dir) > 1 {
		return "", false, ErrTooManyArgs
	} else if len(dir) == 1 {
		dirpath = dir[0]
	}
	header, err := t.Next()
	if err != nil {
		return "", false, err
	}
	path := filepath.Join(dirpath, header.Name)
	filedir, _ := filepath.Split(path)
	if err := checkDir(filedir); err != nil {
		return "", false, err
	}
	if header.FileInfo().IsDir() {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path, true, os.Mkdir(path, 0744)
		}
		return path, true, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return "", false, err
	}
	defer f.Close()
	_, err = io.Copy(f, t)
	return path, false, err
}

// NextBytes advances to the next file name and bytes in the tar archive.
func (t *Reader) NextBytes() (name string, isDir bool, data []byte, err error) {
	header, err := t.Next()
	if err != nil {
		return "", false, nil, err
	}
	name = header.Name
	if header.FileInfo().IsDir() {
		isDir = true
		return
	}
	data, err = ioutil.ReadAll(t)
	return
}

// Close closes the file.
func (t *Reader) Close() error {
	if t.f != nil {
		return t.f.Close()
	}
	return nil
}

func checkDir(dirName string) error {
	if len(dirName) > 0 {
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			return os.MkdirAll(dirName, 0744)
		}
	}
	return nil
}
