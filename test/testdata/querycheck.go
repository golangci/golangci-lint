package testdata

import (
	"database/sql"
)

func main() {
	db, _ := sql.Open("mysql", "test:test@tcp(test:3306)/test")

	test := "a"
	_, err = db.Query("Update * FROM hoge where id = ?", test)           // ERROR "Query not `SELECT` query"
	_, err = db.QueryRow("Update * FROM hoge where id = ?", test)        // ERROR "Query not `SELECT` query"
	_, err = db.QueryContext("Update * FROM hoge where id = ?", test)    // ERROR "Query not `SELECT` query"
	_, err = db.QueryRowContext("Update * FROM hoge where id = ?", test) // ERROR "Query not `SELECT` query"
}
