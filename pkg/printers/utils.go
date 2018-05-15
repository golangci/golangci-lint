package printers

import (
	"os"
	"syscall"
)

var stdOut = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout") // was set to /dev/null
