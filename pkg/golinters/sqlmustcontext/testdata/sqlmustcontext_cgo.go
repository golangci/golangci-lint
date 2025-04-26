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

func _() {
	ctx := context.Background()

	db2, _ := sql.Open("sqlmustcontext", "sqlmustcontext://")

	db2.Ping() // want "use PingContext instead of Ping"
	db2.PingContext(ctx)

	db2.Exec("select * from testdata") // want "use ExecContext instead of Exec"
	db2.ExecContext(ctx, "select * from testdata")

	db2.Prepare("select * from testdata") // want "use PrepareContext instead of Prepare"
	db2.PrepareContext(ctx, "select * from testdata")

	db2.Query("select * from testdata") // want "use QueryContext instead of Query"
	db2.QueryContext(ctx, "select * from testdata")

	db2.QueryRow("select * from testdata") // want "use QueryRowContext instead of QueryRow"
	db2.QueryRowContext(ctx, "select * from testdata")

	// transactions

	tx2, _ := db2.Begin()
	tx2.Exec("select * from testdata") // want "use ExecContext instead of Exec"
	tx2.ExecContext(ctx, "select * from testdata")

	tx2.Prepare("select * from testdata") // want "use PrepareContext instead of Prepare"
	tx2.PrepareContext(ctx, "select * from testdata")

	tx2.Query("select * from testdata") // want "use QueryContext instead of Query"
	tx2.QueryContext(ctx, "select * from testdata")

	tx2.QueryRow("select * from testdata") // want "use QueryRowContext instead of QueryRow"
	tx2.QueryRowContext(ctx, "select * from testdata")

	_ = tx2.Commit()
}
