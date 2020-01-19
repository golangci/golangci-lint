//args: -Erowserrcheck
package testdata

import (
	"database/sql"
)

func RowsErrNotChecked(db *sql.DB) {
	rows, _ := db.Query("select id from tb") // ERROR "rows.Err must be checked"
	for rows.Next() {

	}
}
