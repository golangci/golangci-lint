package codegen

import (
	"bytes"
	"io"
	"os"
)

// ShouldWriteFile returns true of the contents of the file and the given data represent
// and effectively different script. In other words, if the file and the content
// are the same, then do not overwrite the file
func ShouldWriteFile(filename string, content []byte, eq func(a, b io.Reader) bool) bool {
	// check existing file
	fd, err := os.Open(filename)
	if err != nil {
		// maybe doesn't exist --> ok
		// other error --> likely the subsequent write will fail too
		return true
	}
	defer fd.Close()
	return !eq(fd, bytes.NewReader(content))
}
