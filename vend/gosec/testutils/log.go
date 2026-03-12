package testutils

import (
	"bytes"
	"log"
)

// NewLogger returns a logger and the buffer that it will be written to
func NewLogger() (*log.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return log.New(&buf, "", log.Lshortfile), &buf
}
