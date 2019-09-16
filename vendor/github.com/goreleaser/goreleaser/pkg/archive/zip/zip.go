// Package zip implements the Archive interface providing zip archiving
// and compression.
package zip

import (
	"archive/zip"
	"io"
	"os"
)

// Archive zip struct
type Archive struct {
	z *zip.Writer
}

// Close all closeables
func (a Archive) Close() error {
	return a.z.Close()
}

// New zip archive
func New(target io.Writer) Archive {
	return Archive{
		z: zip.NewWriter(target),
	}
}

// Add a file to the zip archive
func (a Archive) Add(name, path string) (err error) {
	file, err := os.Open(path) // #nosec
	if err != nil {
		return
	}
	defer file.Close() // nolint: errcheck
	info, err := file.Stat()
	if err != nil {
		return
	}
	if info.IsDir() {
		return
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = name
	header.Method = zip.Deflate
	w, err := a.z.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, file)
	return err
}
