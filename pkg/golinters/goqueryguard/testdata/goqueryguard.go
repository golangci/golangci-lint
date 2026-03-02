//golangcitest:args -Egoqueryguard
package testdata

import "database/sql"

func queryInLoop(db *sql.DB, ids []int) {
	for range ids {
		_, _ = db.Query("SELECT 1") // want `query-in-loop \[definite\]`
	}
}

func rowsIterationIsNotAQuery(db *sql.DB, ids []int) {
	rows, _ := db.Query("SELECT 1")
	defer rows.Close()

	for range ids {
		_ = rows.Next()
	}
}
