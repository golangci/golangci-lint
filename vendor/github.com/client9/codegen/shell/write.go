package shell

import (
	"github.com/client9/codegen"
)

// ShouldWriteFile returns true of the contents of the file and the given data represent
// and effectively different script. In other words, if the file and the content
// are the same, then do not overwrite the file
func ShouldWriteFile(filename string, content []byte) bool {
	return codegen.ShouldWriteFile(filename, content, Equal)
}
