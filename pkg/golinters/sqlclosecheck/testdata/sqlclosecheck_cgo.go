//golangcitest:args -Esqlclosecheck
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
	"context"
	"database/sql"
	"log"
	"strings"
	"unsafe"
)

var (
	ctx    context.Context
	db     *sql.DB
	age    = 27
	userID = 43
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	rows, err := db.QueryContext(ctx, "SELECT name FROM users WHERE age=?", age) // want "Rows/Stmt/NamedStmt was not closed"
	if err != nil {
		log.Fatal(err)
	}
	// defer rows.Close()

	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s are %d years old", strings.Join(names, ", "), age)
}
