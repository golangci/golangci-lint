// Package archive provides tar.gz and zip archiving
package archive

import (
	"os"
	"strings"

	"github.com/goreleaser/goreleaser/pkg/archive/gzip"
	"github.com/goreleaser/goreleaser/pkg/archive/targz"
	"github.com/goreleaser/goreleaser/pkg/archive/zip"
)

// Archive represents a compression archive files from disk can be written to.
type Archive interface {
	Close() error
	Add(name, path string) error
}

// New archive.
func New(file *os.File) Archive {
	if strings.HasSuffix(file.Name(), ".tar.gz") {
		return targz.New(file)
	}
	if strings.HasSuffix(file.Name(), ".gz") {
		return gzip.New(file)
	}
	if strings.HasSuffix(file.Name(), ".zip") {
		return zip.New(file)
	}
	return targz.New(file)
}
