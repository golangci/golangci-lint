//golangcitest:args -Ewastedassign
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
	"strings"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func pa(x int) int {
	return x + 1
}

func multiple(val interface{}, times uint) interface{} {

	switch hogehoge := val.(type) {
	case int:
		return 12
	case string:
		return strings.Repeat(hogehoge, int(times))
	default:
		return nil
	}
}

func noUseParams(params string) int {
	a := 12
	println(a)
	return a
}

func f(param int) int {
	println(param)
	useOutOfIf := 1212121 // want "assigned to useOutOfIf, but reassigned without using the value"
	ret := 0
	if false {
		useOutOfIf = 200 // want "assigned to useOutOfIf, but never used afterwards"
		return 0
	} else if param == 100 {
		useOutOfIf = 100 // want "assigned to useOutOfIf, but reassigned without using the value"
		useOutOfIf = 201
		useOutOfIf = pa(useOutOfIf)
		useOutOfIf += 200 // want "assigned to useOutOfIf, but reassigned without using the value"
	} else {
		useOutOfIf = 100
		useOutOfIf += 100
		useOutOfIf = pa(useOutOfIf)
		useOutOfIf += 200 // want "assigned to useOutOfIf, but reassigned without using the value"
	}

	if false {
		useOutOfIf = 200 // want "assigned to useOutOfIf, but never used afterwards"
		return 0
	} else if param == 200 {
		useOutOfIf = 100 // want "assigned to useOutOfIf, but reassigned without using the value"
		useOutOfIf = 201
		useOutOfIf = pa(useOutOfIf)
		useOutOfIf += 200
	} else {
		useOutOfIf = 100
		useOutOfIf += 100
		useOutOfIf = pa(useOutOfIf)
		useOutOfIf += 200
	}
	// useOutOfIf = 12
	println(useOutOfIf)
	useOutOfIf = 192
	useOutOfIf += 100
	useOutOfIf += 200 // want "assigned to useOutOfIf, but never used afterwards"
	return ret
}

func checkLoopTest() int {
	hoge := 12
	noUse := 1111
	println(noUse)

	noUse = 1111 // want "assigned to noUse, but never used afterwards"
	for {
		if hoge == 14 {
			break
		}
		hoge = hoge + 1
	}
	return hoge
}

func r(param int) int {
	println(param)
	useOutOfIf := 1212121
	ret := 0
	if false {
		useOutOfIf = 200 // want "assigned to useOutOfIf, but never used afterwards"
		return 0
	} else if param == 100 {
		ret = useOutOfIf
	} else if param == 200 {
		useOutOfIf = 100 // want "assigned to useOutOfIf, but reassigned without using the value"
		useOutOfIf = 100
		useOutOfIf = pa(useOutOfIf)
		useOutOfIf += 200 // want "assigned to useOutOfIf, but reassigned without using the value"
	}
	useOutOfIf = 12
	println(useOutOfIf)
	useOutOfIf = 192
	useOutOfIf += 100
	useOutOfIf += 200 // want "assigned to useOutOfIf, but never used afterwards"
	return ret
}

func mugen() {
	var i int
	var hoge int
	for {
		hoge = 5 // want "assigned to hoge, but reassigned without using the value"
		// break
	}

	println(i)
	println(hoge)
	return
}

func mugenG[T ~int](hoge T) {
	var i int
	for {
		hoge = 5 // want "assigned to hoge, but reassigned without using the value"
		// break
	}

	println(i)
	println(hoge)
	return
}

func noMugen() {
	var i int
	var hoge int
	for {
		hoge = 5
		break
	}

	println(i)
	println(hoge)
	return
}

func reassignInsideLoop() {
	bar := func(b []byte) ([]byte, error) { return b, nil }
	var err error
	var rest []byte
	for {
		rest, err = bar(rest)
		if err == nil {
			break
		}
	}
	return
}

func reassignInsideLoop2() {
	var x int = 0
	var y int = 1
	for i := 1; i < 3; i++ {
		x += y
		y *= 2 * i
	}
	println(x)
}
