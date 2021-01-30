//args: -Erowserrcheck
package testdata

import (
	"database/sql"
	"fmt"
	"math/rand"
)

func RowsErrNotChecked(db *sql.DB) {
	rows, _ := db.Query("select id from tb") // ERROR "rows.Err must be checked"
	for rows.Next() {

	}
}

func issue943_1() {
	db, err := sql.Open("postgres", "postgres://localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
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

func issue943_1() {
	db, err := sql.Open("postgres", "postgres://localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var rows *sql.Rows
	if rand.Float64() < 0.5 {
		rows, err = db.Query("select 1") // ERROR "rows.Err must be checked"
	} else {
		rows, err = db.Query("select 2") // ERROR "rows.Err must be checked"
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
