// args: -Eexecinquery
package testdata

import (
	"context"
	"database/sql"
)

func execInQuery(db *sql.DB) {
	test := "a"

	_, err := db.Query("Update * FROM hoge where id = ?", test) // ERROR "It's better to use Execute method instead of Query method to execute `UPDATE` query"
	if err != nil {
		return
	}

	db.QueryRow("Update * FROM hoge where id = ?", test) // ERROR "It's better to use Execute method instead of QueryRow method to execute `UPDATE` query"
	if err != nil {
		return
	}

	ctx := context.Background()

	_, err = db.QueryContext(ctx, "Update * FROM hoge where id = ?", test) // ERROR "It's better to use Execute method instead of QueryContext method to execute `UPDATE` query "
	if err != nil {
		return
	}

	db.QueryRowContext(ctx, "Update * FROM hoge where id = ?", test) // ERROR "It's better to use Execute method instead of QueryRowContext method to execute `UPDATE` query"
}
