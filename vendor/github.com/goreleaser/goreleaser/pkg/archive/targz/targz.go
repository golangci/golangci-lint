// Package targz implements the Archive interface providing tar.gz archiving
// and compression.
package targz

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

// Archive as tar.gz
type Archive struct {
	gw *gzip.Writer
	tw *tar.Writer
}

// Close all closeables
func (a Archive) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	return a.gw.Close()
}

// New tar.gz archive
func New(target io.Writer) Archive {
	gw := gzip.NewWriter(target)
	tw := tar.NewWriter(gw)
	return Archive{
		gw: gw,
		tw: tw,
	}
}

// Add file to the archive
func (a Archive) Add(name, path string) error {
	file, err := os.Open(path) // #nosec
	if err != nil {
		return err
	}
	defer file.Close() // nolint: errcheck
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, name)
	if err != nil {
		return err
	}
	header.Name = name
	if err = a.tw.WriteHeader(header); err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	_, err = io.Copy(a.tw, file)
	return err
}
