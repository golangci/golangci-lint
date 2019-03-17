package fsutils

import (
	"fmt"
	"io/ioutil"

	"github.com/golangci/golangci-lint/pkg/logutils"

	"github.com/pkg/errors"
)

type FileCache struct {
	files map[string][]byte
}

func NewFileCache() *FileCache {
	return &FileCache{
		files: map[string][]byte{},
	}
}

func (fc *FileCache) GetFileBytes(filePath string) ([]byte, error) {
	cachedBytes := fc.files[filePath]
	if cachedBytes != nil {
		return cachedBytes, nil
	}

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read file %s", filePath)
	}

	fc.files[filePath] = fileBytes
	return fileBytes, nil
}

func prettifyBytesCount(n int) string {
	const (
		Multiplexer = 1024
		KiB         = 1 * Multiplexer
		MiB         = KiB * Multiplexer
		GiB         = MiB * Multiplexer
	)

	if n >= GiB {
		return fmt.Sprintf("%.1fGiB", float64(n)/GiB)
	}
	if n >= MiB {
		return fmt.Sprintf("%.1fMiB", float64(n)/MiB)
	}
	if n >= KiB {
		return fmt.Sprintf("%.1fKiB", float64(n)/KiB)
	}
	return fmt.Sprintf("%dB", n)
}

func (fc *FileCache) PrintStats(log logutils.Log) {
	var size int
	for _, fileBytes := range fc.files {
		size += len(fileBytes)
	}

	log.Infof("File cache stats: %d entries of total size %s", len(fc.files), prettifyBytesCount(size))
}
