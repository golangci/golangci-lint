package printers

import (
	"os"
	"syscall"
)

var StdOut = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout") // was set to /dev/null
