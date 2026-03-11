package testutils

import "github.com/securego/gosec/v2"

// SampleCodeCgo - Cgo file sample
var SampleCodeCgo = []CodeSample{
	{[]string{`
package main

import (
        "fmt"
        "unsafe"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

int printData(unsigned char *data) {
    return printf("cData: %lu \"%s\"\n", (long unsigned int)strlen(data), data);
}
*/
import "C"

func main() {
        // Allocate C data buffer.
        width, height := 8, 2
        lenData := width * height
        // add string terminating null byte
        cData := (*C.uchar)(C.calloc(C.size_t(lenData+1), C.sizeof_uchar))

        // When no longer in use, free C allocations.
        defer C.free(unsafe.Pointer(cData))

        // Go slice reference to C data buffer,
        // minus string terminating null byte
        gData := (*[1 << 30]byte)(unsafe.Pointer(cData))[:lenData:lenData]

        // Write and read cData via gData.
        for i := range gData {
                gData[i] = '.'
        }
        copy(gData[0:], "Data")
        gData[len(gData)-1] = 'X'
        fmt.Printf("gData: %d %q\n", len(gData), gData)
        C.printData(cData)
}
`}, 0, gosec.NewConfig()},
}
