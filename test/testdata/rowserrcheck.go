//args: -Ebodyclose
package testdata

import (
	"database/sql"
)

func RowsErrNotChecked(db *sql.DB) {
	rows, _ := db.Query("select id from tb") // ERROR "response body must be closed"
	for rows.Next() {

	}
}
