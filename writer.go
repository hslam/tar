// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
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
		info, err := os.Stat(name)
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := w.TarDir(name); err != nil {
				return err
			}
			continue
		}
		if err := w.tarFile(name, info); err != nil {
			return err
		}
	}
	return nil
}

// TarDir tars the dir.
func (w *Writer) TarDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.tarDir(path, info)
		}
		return w.tarFile(path, info)
	})
}

// TarFile writes a file to the tar file.
func (w *Writer) TarFile(name string) error {
	info, err := os.Stat(name)
	if err != nil {
		return err
	}
	return w.tarFile(name, info)
}

func (w *Writer) tarDir(path string, info os.FileInfo) error {
	hdr, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	hdr.Name = path
	return w.WriteHeader(hdr)
}

func (w *Writer) tarFile(path string, info os.FileInfo) error {
	hdr, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	hdr.Name = path
	err = w.WriteHeader(hdr)
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
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
