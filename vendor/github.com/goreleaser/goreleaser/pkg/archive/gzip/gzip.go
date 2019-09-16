// Package gzip implements the Archive interface providing gz archiving
// and compression.
package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// Archive as gz
type Archive struct {
	gw *gzip.Writer
}

// Close all closeables
func (a Archive) Close() error {
	return a.gw.Close()
}

// New gz archive
func New(target io.Writer) Archive {
	gw := gzip.NewWriter(target)
	return Archive{
		gw: gw,
	}
}

// Add file to the archive
func (a Archive) Add(name, path string) error {
	if a.gw.Header.Name != "" {
		return fmt.Errorf("gzip: failed to add %s, only one file can be archived in gz format", name)
	}
	file, err := os.Open(path) // #nosec
	if err != nil {
		return err
	}
	defer file.Close() // nolint: errcheck
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	a.gw.Header.Name = name
	a.gw.Header.ModTime = info.ModTime()
	_, err = io.Copy(a.gw, file)
	return err
}
