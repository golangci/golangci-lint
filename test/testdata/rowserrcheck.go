//golangcitest:args -Erowserrcheck
package testdata

import (
	"database/sql"
	"fmt"
	"math/rand"
)

func RowsErrNotChecked(db *sql.DB) {
	rows, _ := db.Query("select id from tb") // want "rows.Err must be checked"
	for rows.Next() {

	}
}

func issue943(db *sql.DB) {
	var rows *sql.Rows
	var err error

	if rand.Float64() < 0.5 {
		rows, err = db.Query("select 1")
	} else {
		rows, err = db.Query("select 2")
	}
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		fmt.Println("new rows")
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
