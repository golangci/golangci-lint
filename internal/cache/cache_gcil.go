package cache

import (
	"errors"

	"github.com/golangci/golangci-lint/internal/robustio"
)

// IsErrMissing allows to access to the internal error.
// TODO(ldez) the handling of this error inside runner_action.go should be refactored.
func IsErrMissing(err error) bool {
	var errENF *entryNotFoundError
	return errors.As(err, &errENF)
}

func (c *Cache) readFileCGIL(outputFile string, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	return robustio.ReadFile(outputFile)
}
