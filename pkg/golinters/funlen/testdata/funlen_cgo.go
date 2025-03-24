//golangcitest:args -Efunlen
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

func TooManyLines() { // want `Function 'TooManyLines' is too long \(70 > 60\)`
	t := struct {
		A0  string
		A1  string
		A2  string
		A3  string
		A4  string
		A5  string
		A6  string
		A7  string
		A8  string
		A9  string
		A10 string
		A11 string
		A12 string
		A13 string
		A14 string
		A15 string
		A16 string
		A17 string
		A18 string
		A19 string
		A20 string
		A21 string
		A22 string
		A23 string
		A24 string
		A25 string
		A26 string
		A27 string
		A28 string
		A29 string
		A30 string
		A31 string
		A32 string
	}{
		A0:  "a",
		A1:  "a",
		A2:  "a",
		A3:  "a",
		A4:  "a",
		A5:  "a",
		A6:  "a",
		A7:  "a",
		A8:  "a",
		A9:  "a",
		A10: "a",
		A11: "a",
		A12: "a",
		A13: "a",
		A14: "a",
		A15: "a",
		A16: "a",
		A17: "a",
		A18: "a",
		A19: "a",
		A20: "a",
		A21: "a",
		A22: "a",
		A23: "a",
		A24: "a",
		A25: "a",
		A26: "a",
		A27: "a",
		A28: "a",
		A29: "a",
		A30: "a",
		A31: "a",
		A32: "a",
	}
	_ = t
}

func TooManyStatements() { // want `Function 'TooManyStatements' has too many statements \(46 > 40\)`
	a0 := 1
	a1 := 1
	a2 := 1
	a3 := 1
	a4 := 1
	a5 := 1
	a6 := 1
	a7 := 1
	a8 := 1
	a9 := 1
	a10 := 1
	a11 := 1
	a12 := 1
	a13 := 1
	a14 := 1
	a15 := 1
	a16 := 1
	a17 := 1
	a18 := 1
	a19 := 1
	a20 := 1
	a21 := 1
	a22 := 1
	_ = a0
	_ = a1
	_ = a2
	_ = a3
	_ = a4
	_ = a5
	_ = a6
	_ = a7
	_ = a8
	_ = a9
	_ = a10
	_ = a11
	_ = a12
	_ = a13
	_ = a14
	_ = a15
	_ = a16
	_ = a17
	_ = a18
	_ = a19
	_ = a20
	_ = a21
	_ = a22
}

func withComments() {
	// Comment 1
	// Comment 2
	// Comment 3
	// Comment 4
	// Comment 5
	// Comment 6
	// Comment 7
	// Comment 8
	// Comment 9
	// Comment 10
	// Comment 11
	// Comment 12
	// Comment 13
	// Comment 14
	// Comment 15
	// Comment 16
	// Comment 17
	// Comment 18
	// Comment 19
	// Comment 20
	// Comment 21
	// Comment 22
	// Comment 23
	// Comment 24
	// Comment 25
	// Comment 26
	// Comment 27
	// Comment 28
	// Comment 29
	// Comment 30
	// Comment 31
	// Comment 32
	// Comment 33
	// Comment 34
	// Comment 35
	// Comment 36
	// Comment 37
	// Comment 38
	// Comment 39
	// Comment 40
	// Comment 41
	// Comment 42
	// Comment 43
	// Comment 44
	// Comment 45
	// Comment 46
	// Comment 47
	// Comment 48
	// Comment 49
	// Comment 50
	// Comment 51
	// Comment 52
	// Comment 53
	// Comment 54
	// Comment 55
	// Comment 56
	// Comment 57
	// Comment 58
	// Comment 59
	// Comment 60
	print("Hello, world!")
}
