//golangcitest:args -Egosmopolitan
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
	"fmt"
	"time"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type col struct {
	// struct tag should not get reported
	Foo string `gorm:"column:bar;not null;comment:'不应该报告这一行'"`
}

func _() {
	fmt.Println("hello world")
	fmt.Println("你好，世界") // want `string literal contains rune in Han script`
	fmt.Println("こんにちは、セカイ")

	_ = col{Foo: "hello"}
	_ = col{Foo: "你好"} // want `string literal contains rune in Han script`

	x := time.Local // want `usage of time.Local`
	_ = time.Now().In(x)
	_ = time.Date(2023, 1, 2, 3, 4, 5, 678901234, time.Local) // want `usage of time.Local`
}
