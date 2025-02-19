// Package pkg ...
package pkg

import (
	"log"
	"unsafe"
)

// F ...
func F() {
	x := 123 // of type int
	p := unsafe.Pointer(&x)
	pp := &p // of type *unsafe.Pointer
	p = unsafe.Pointer(pp)
	log.Print(p)
}
