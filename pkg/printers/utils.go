package printers

import (
	"io"
	"os"
	"syscall"
)

func getOutWriter() io.Writer {
	return os.NewFile(uintptr(syscall.Stdout), "/dev/stdout") // was set to /dev/null
}
