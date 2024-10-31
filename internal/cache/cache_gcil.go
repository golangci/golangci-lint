package cache

import (
	"errors"

	"github.com/golangci/golangci-lint/internal/robustio"
)

func IsErrMissing(err error) bool {
	return errors.Is(err, errMissing)
}

func (c *Cache) readFileCGIL(outputFile string, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	return robustio.ReadFile(outputFile)
}
