//args: -Eexecinquery
package testdata

import (
	"database/sql"

	"golang.org/x/net/context"
)

func execinquery() {
	db, _ := sql.Open("mysql", "test:test@tcp(test:3306)/test")

	test := "a"

	_, err := db.Query("Update * FROM hoge where id = ?", test) // ERROR "Query not `SELECT` query"
	if err != nil {
		return
	}

	db.QueryRow("Update * FROM hoge where id = ?", test) // ERROR "Query not `SELECT` query"
	if err != nil {
		return
	}

	ctx := context.Background()

	_, err = db.QueryContext(ctx, "Update * FROM hoge where id = ?", test) // ERROR "Query not `SELECT` query"
	if err != nil {
		return
	}

	db.QueryRowContext(ctx, "Update * FROM hoge where id = ?", test) // ERROR "Query not `SELECT` query"

}
