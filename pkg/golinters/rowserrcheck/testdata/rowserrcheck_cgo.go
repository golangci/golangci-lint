//golangcitest:args -Erowserrcheck
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
	"database/sql"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _(db *sql.DB) {
	rows, _ := db.Query("select id from tb") // want "rows.Err must be checked"
	for rows.Next() {

	}
}
