// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package tar implements access to tar archives.
package tar

// TarSuffix represents the tar suffix.
const TarSuffix = ".tar"

// TargzSuffix represents the tar gzip suffix.
const TargzSuffix = ".tar.gz"

// Tar tars all the paths to the tar file.
func Tar(tar string, paths ...string) error {
	tw, err := NewFileWriter(tar)
	if err != nil {
		return err
	}
	defer tw.Close()
	return tw.Tar(paths...)
}

// Untar untars all the files to dir.
func Untar(tar string, dir ...string) ([]string, []string, error) {
	tr, err := NewFileReader(tar)
	if err != nil {
		return nil, nil, err
	}
	defer tr.Close()
	return tr.Untar(dir...)
}

// Targz tars all the paths to the tar gzip file.
func Targz(targz string, paths ...string) error {
	tw, err := NewGzipFileWriter(targz)
	if err != nil {
		return err
	}
	defer tw.Close()
	return tw.Tar(paths...)
}

// Untargz untars all the files to dir.
func Untargz(targz string, dir ...string) ([]string, []string, error) {
	tr, err := NewGzipFileReader(targz)
	if err != nil {
		return nil, nil, err
	}
	defer tr.Close()
	return tr.Untar(dir...)
}
