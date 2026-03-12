package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG103 find instances of unsafe blocks for auditing purposes
var SampleCodeG103 = []CodeSample{
	{[]string{`
package main

import (
	"fmt"
	"unsafe"
)

type Fake struct{}

func (Fake) Good() {}

func main() {
	unsafeM := Fake{}
	unsafeM.Good()
	intArray := [...]int{1, 2}
	fmt.Printf("\nintArray: %v\n", intArray)
	intPtr := &intArray[0]
	fmt.Printf("\nintPtr=%p, *intPtr=%d.\n", intPtr, *intPtr)
	addressHolder := uintptr(unsafe.Pointer(intPtr)) 
	intPtr = (*int)(unsafe.Pointer(addressHolder))
	fmt.Printf("\nintPtr=%p, *intPtr=%d.\n\n", intPtr, *intPtr)
}
`}, 2, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	chars := [...]byte{1, 2}
	charsPtr := &chars[0]
	str := unsafe.String(charsPtr, len(chars))
	fmt.Printf("%s\n", str)
	ptr := unsafe.StringData(str)
	fmt.Printf("ptr: %p\n", ptr)
}
`}, 2, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	chars := [...]byte{1, 2}
	charsPtr := &chars[0]
	slice := unsafe.Slice(charsPtr, len(chars))
	fmt.Printf("%v\n", slice)
	ptr := unsafe.SliceData(slice)
	fmt.Printf("ptr: %p\n", ptr)
}
`}, 2, gosec.NewConfig()},
}
