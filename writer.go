// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Writer provides sequential writing of a tar archive.
type Writer struct {
	*tar.Writer
	gw *gzip.Writer
	f  *os.File
}

// NewWriter creates a new Writer writing to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{Writer: tar.NewWriter(w)}
}

// NewGzipWriter creates a new gzip Writer writing to w.
func NewGzipWriter(w io.Writer) *Writer {
	gw := gzip.NewWriter(w)
	return &Writer{Writer: tar.NewWriter(gw), gw: gw}
}

// NewFileWriter creates a new Writer writing to file.
func NewFileWriter(name string) (*Writer, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &Writer{Writer: tar.NewWriter(f), f: f}, nil
}

// NewGzipFileWriter creates a new gzip Writer writing to file.
func NewGzipFileWriter(name string) (*Writer, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	gw := gzip.NewWriter(f)
	return &Writer{Writer: tar.NewWriter(gw), gw: gw, f: f}, nil
}

// Flush finishes writing the current file's block padding.
// The current file must be fully written before Flush can be called.
//
// This is unnecessary as the next call to WriteHeader or Close
// will implicitly flush out the file's padding.
func (w *Writer) Flush() error {
	err := w.Writer.Flush()
	if err != nil {
		return err
	}
	if w.gw != nil {
		return w.gw.Flush()
	}
	return nil
}

// Tar tars all the paths to the tar file.
func (w *Writer) Tar(paths ...string) error {
	for _, name := range paths {
		if err := w.tarPath(name); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) tarPath(name string) error {
	info, err := os.Stat(name)
	if err == nil {
		if info.IsDir() {
			err = w.TarDir(name)
		} else {
			err = w.tarFile(name, "", info)
		}
	}
	return err
}

// TarDir tars a dir to the tar file.
func (w *Writer) TarDir(dir string) error {
	var base = filepath.Base(dir)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			var rel string
			rel, err = filepath.Rel(dir, path)
			if err == nil {
				if info.IsDir() {
					return w.tarDir(path, filepath.Join(base, rel), info)
				}
				err = w.tarFile(path, filepath.Join(base, rel), info)
			}
		}
		return err
	})
}

// TarFile writes a file to the tar file.
func (w *Writer) TarFile(name string) error {
	info, err := os.Stat(name)
	if err != nil {
		return err
	}
	return w.tarFile(name, "", info)
}

func (w *Writer) tarDir(path, name string, info os.FileInfo) (err error) {
	var hdr *tar.Header
	hdr, err = tar.FileInfoHeader(info, "")
	if err == nil {
		if len(name) > 0 {
			hdr.Name = name
		}
		err = w.WriteHeader(hdr)
	}
	return err
}

func (w *Writer) tarFile(path, name string, info os.FileInfo) (err error) {
	var hdr *tar.Header
	hdr, err = tar.FileInfoHeader(info, "")
	if err == nil {
		if len(name) > 0 {
			hdr.Name = name
		}
		err = w.WriteHeader(hdr)
		if err == nil {
			var f *os.File
			f, err = os.Open(path)
			if err == nil {
				_, err = io.Copy(w, f)
				f.Close()
			}
		}
	}
	return err
}

// TarBytes tars a file with the file name and data.
func (w *Writer) TarBytes(name string, data []byte) (err error) {
	hdr := &tar.Header{
		Name:    name,
		Size:    int64(len(data)),
		Mode:    0666,
		ModTime: time.Now(),
	}
	err = w.WriteHeader(hdr)
	if err == nil {
		_, err = w.Write(data)
	}
	return err
}

// Close closes the tar archive by flushing the padding, and writing the footer.
// If the current file (from a prior call to WriteHeader) is not fully written,
// then this returns an error.
func (w *Writer) Close() error {
	err := w.Writer.Close()
	if w.gw != nil {
		w.gw.Close()
	}
	if w.f != nil {
		w.f.Close()
	}
	return err
}
