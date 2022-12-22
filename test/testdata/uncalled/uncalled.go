//golangcitest:args -Euncalled
package testdata

import (
	"database/sql"
)

func RowsErrNotChecked(db *sql.DB) {
	rows, err := db.Query("select id from tb") // want "rows.Err\\(\\) must be called"
	if err != nil {
		// Handle error.
	}

	for rows.Next() {
		// Handle row.
	}
}
