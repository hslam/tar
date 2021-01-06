// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package tar

import (
	"os"
	"testing"
)

func TestTar(t *testing.T) {
	name := "file"
	tname := name + TarSuffix
	defer os.Remove(name)
	defer os.Remove(tname)
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
	}
	contents := "Hello World"
	file.Write([]byte(contents))
	file.Close()
	Tar(tname, name)
	os.Remove(name)
	Untar(tname)
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, len(contents))
	f.Read(buf)
	if string(buf) != contents {
		t.Error(string(buf), contents)
	}
}

func TestTargz(t *testing.T) {
	name := "file"
	tname := name + TargzSuffix
	defer os.Remove(name)
	defer os.Remove(tname)
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
	}
	contents := "Hello World"
	file.Write([]byte(contents))
	file.Close()
	Targz(tname, name)
	os.Remove(name)
	Untargz(tname)
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, len(contents))
	f.Read(buf)
	if string(buf) != contents {
		t.Error(string(buf), contents)
	}
}

func TestReadWriter(t *testing.T) {
	name := "file"
	tname := name + TarSuffix
	defer os.Remove(name)
	defer os.Remove(tname)
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
	}
	contents := "Hello World"
	file.Write([]byte(contents))
	file.Close()
	w, err := os.Create(tname)
	if err != nil {
		t.Error(err)
	}
	tw := NewWriter(w)
	tw.TarFile(name)
	tw.Flush()
	tw.Close()
	w.Close()
	os.Remove(name)

	r, err := os.Open(tname)
	if err != nil {
		t.Error(err)
	}
	tr := NewReader(r)
	tr.NextFile()
	tr.Close()
	r.Close()

	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, len(contents))
	f.Read(buf)
	if string(buf) != contents {
		t.Error(string(buf), contents)
	}
}

func TestGzipReadWriter(t *testing.T) {
	name := "file"
	tname := name + TargzSuffix
	defer os.Remove(name)
	defer os.Remove(tname)
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
	}
	contents := "Hello World"
	file.Write([]byte(contents))
	file.Close()
	w, err := os.Create(tname)
	if err != nil {
		t.Error(err)
	}
	tw := NewGzipWriter(w)
	tw.TarFile(name)
	tw.Flush()
	tw.Close()
	w.Close()
	os.Remove(name)

	r, err := os.Open(tname)
	if err != nil {
		t.Error(err)
	}
	tr := NewGzipReader(r)
	tr.NextFile()
	tr.Close()
	r.Close()

	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, len(contents))
	f.Read(buf)
	if string(buf) != contents {
		t.Error(string(buf), contents)
	}
}

func TestTarBytes(t *testing.T) {
	name := "file"
	name1 := "file1"
	tname := "file.tar"
	defer os.Remove(tname)
	contents := "Hello World"
	contents1 := "Hello World1"

	w, err := os.Create(tname)
	if err != nil {
		t.Error(err)
	}
	tw := NewWriter(w)
	tw.TarBytes(name, []byte(contents))
	tw.TarBytes(name1, []byte(contents1))
	tw.Flush()
	tw.Close()
	w.Close()
	r, err := os.Open(tname)
	if err != nil {
		t.Error(err)
	}
	tr := NewReader(r)
	n, data, err := tr.NextBytes()
	if err != nil {
		t.Error(err)
	}
	if n != name {
		t.Error(n, name)
	}
	if string(data) != contents {
		t.Error(string(data), contents)
	}
	n1, data1, err := tr.NextBytes()
	if err != nil {
		t.Error(err)
	}
	if n1 != name1 {
		t.Error(n1, name1)
	}
	if string(data1) != contents1 {
		t.Error(string(data1), contents1)
	}
	_, _, err = tr.NextBytes()
	if err == nil {
		t.Error()
	}
	tr.Close()
	r.Close()
}

func TestTarDir(t *testing.T) {
	dir := "dir"
	err := checkDir(dir)
	if err != nil {
		t.Error(err)
	}
	name := dir + "/" + "file"
	name1 := dir + "/" + "file1"
	tname := "file" + TarSuffix
	defer os.RemoveAll(dir)
	defer os.Remove(tname)
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
	}
	contents := "Hello World"
	file.Write([]byte(contents))
	file.Close()
	file1, err := os.Create(name1)
	if err != nil {
		t.Error(err)
	}
	file1.Write([]byte(contents))
	file1.Close()
	w, err := os.Create(tname)
	if err != nil {
		t.Error(err)
	}
	tw := NewWriter(w)
	tw.TarDir(dir)
	tw.Flush()
	tw.Close()
	w.Close()
	os.RemoveAll(dir)
	r, err := os.Open(tname)
	if err != nil {
		t.Error(err)
	}
	tr := NewReader(r)
	tr.Untar()
	tr.Close()
	r.Close()

	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, len(contents))
	f.Read(buf)
	if string(buf) != contents {
		t.Error(string(buf), contents)
	}

	f1, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf = make([]byte, len(contents))
	f1.Read(buf)
	if string(buf) != contents {
		t.Error(string(buf), contents)
	}
}

func TestTarPaths(t *testing.T) {
	dir := "dir"
	err := checkDir(dir)
	if err != nil {
		t.Error(err)
	}
	name := dir + "/" + "file"
	name1 := "file1"
	tname := "file" + TarSuffix
	defer os.RemoveAll(dir)
	defer os.Remove(name1)
	defer os.Remove(tname)
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
	}
	contents := "Hello World"
	file.Write([]byte(contents))
	file.Close()
	file1, err := os.Create(name1)
	if err != nil {
		t.Error(err)
	}
	file1.Write([]byte(contents))
	file1.Close()
	w, err := os.Create(tname)
	if err != nil {
		t.Error(err)
	}
	tw := NewWriter(w)
	tw.Tar(dir, name1)
	tw.Flush()
	tw.Close()
	w.Close()
	os.RemoveAll(dir)
	os.Remove(name1)
	r, err := os.Open(tname)
	if err != nil {
		t.Error(err)
	}
	tr := NewReader(r)
	tr.Untar()
	tr.Close()
	r.Close()

	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, len(contents))
	f.Read(buf)
	if string(buf) != contents {
		t.Error(string(buf), contents)
	}

	f1, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf = make([]byte, len(contents))
	f1.Read(buf)
	if string(buf) != contents {
		t.Error(string(buf), contents)
	}
}
