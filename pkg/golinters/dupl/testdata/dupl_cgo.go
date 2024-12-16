//golangcitest:args -Edupl
//golangcitest:config_path testdata/dupl.yml
package testdata

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import (
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type DuplLogger struct{}

func (DuplLogger) level() int {
	return 1
}

func (DuplLogger) Debug(args ...interface{}) {}
func (DuplLogger) Info(args ...interface{})  {}

func (logger *DuplLogger) First(args ...interface{}) { // want "34-43 lines are duplicate of `.*dupl_cgo.go:45-54`"
	if logger.level() >= 0 {
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
		logger.Debug(args...)
	}
}

func (logger *DuplLogger) Second(args ...interface{}) { // want "45-54 lines are duplicate of `.*dupl_cgo.go:34-43`"
	if logger.level() >= 1 {
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
		logger.Info(args...)
	}
}
