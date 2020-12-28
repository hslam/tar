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
	contexts := "Hello World"
	file.Write([]byte(contexts))
	file.Close()
	Tar(tname, name)
	os.Remove(name)
	Untar(tname)
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, len(contexts))
	f.Read(buf)
	if string(buf) != contexts {
		t.Error(string(buf), contexts)
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
	contexts := "Hello World"
	file.Write([]byte(contexts))
	file.Close()
	Targz(tname, name)
	os.Remove(name)
	Untargz(tname)
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, len(contexts))
	f.Read(buf)
	if string(buf) != contexts {
		t.Error(string(buf), contexts)
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
	contexts := "Hello World"
	file.Write([]byte(contexts))
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
	buf := make([]byte, len(contexts))
	f.Read(buf)
	if string(buf) != contexts {
		t.Error(string(buf), contexts)
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
	contexts := "Hello World"
	file.Write([]byte(contexts))
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
	buf := make([]byte, len(contexts))
	f.Read(buf)
	if string(buf) != contexts {
		t.Error(string(buf), contexts)
	}
}

func TestTarDir(t *testing.T) {
	dir := "dir"
	err := checkDir(dir)
	if err != nil {
		t.Error(err)
	}
	name := dir + "/" + "file"
	tname := "file" + TarSuffix
	defer os.RemoveAll(dir)
	defer os.Remove(tname)
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
	}
	contexts := "Hello World"
	file.Write([]byte(contexts))
	file.Close()
	w, err := os.Create(tname)
	if err != nil {
		t.Error(err)
	}
	tw := NewWriter(w)
	tw.TarDir(dir)
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
	buf := make([]byte, len(contexts))
	f.Read(buf)
	if string(buf) != contexts {
		t.Error(string(buf), contexts)
	}
}
