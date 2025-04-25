//golangcitest:args -Erecvcheck
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
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	ctx := context.Background()

	db, _ := sql.Open("sqlmustcontext", "sqlmustcontext://")

	db.Ping() // want "use PingContext instead of Ping"
	db.PingContext(ctx)

	db.Exec("select * from testdata") // want "use ExecContext instead of Exec"
	db.ExecContext(ctx, "select * from testdata")

	db.Prepare("select * from testdata") // want "use PrepareContext instead of Prepare"
	db.PrepareContext(ctx, "select * from testdata")

	db.Query("select * from testdata") // want "use QueryContext instead of Query"
	db.QueryContext(ctx, "select * from testdata")

	db.QueryRow("select * from testdata") // want "use QueryRowContext instead of QueryRow"
	db.QueryRowContext(ctx, "select * from testdata")

	// transactions

	tx, _ := db.Begin()
	tx.Exec("select * from testdata") // want "use ExecContext instead of Exec"
	tx.ExecContext(ctx, "select * from testdata")

	tx.Prepare("select * from testdata") // want "use PrepareContext instead of Prepare"
	tx.PrepareContext(ctx, "select * from testdata")

	tx.Query("select * from testdata") // want "use QueryContext instead of Query"
	tx.QueryContext(ctx, "select * from testdata")

	tx.QueryRow("select * from testdata") // want "use QueryRowContext instead of QueryRow"
	tx.QueryRowContext(ctx, "select * from testdata")

	_ = tx.Commit()
}
