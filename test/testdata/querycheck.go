package testdata

import (
	"database/sql"
)

func main() {
	db, _ := sql.Open("mysql", "test:test@tcp(test:3306)/test")

	test := "a"
	_, err = db.Query("Update * FROM hoge where id = ?", test) // ERROR "Query not `SELECT` query"
}
