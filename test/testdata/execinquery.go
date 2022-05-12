// args: -Eexecinquery
package testdata

import (
	"context"
	"database/sql"
)

func execInQuery(db *sql.DB) {
	test := "a"

	_, err := db.Query("Update * FROM hoge where id = ?", test) // ERROR "Use Exec instead of Query to execute `UPDATE` query"
	if err != nil {
		return
	}

	db.QueryRow("Update * FROM hoge where id = ?", test) // ERROR "Use Exec instead of QueryRow to execute `UPDATE` query"
	if err != nil {
		return
	}

	ctx := context.Background()

	_, err = db.QueryContext(ctx, "Update * FROM hoge where id = ?", test) // ERROR "Use ExecContext instead of QueryContext to execute `UPDATE` query"
	if err != nil {
		return
	}

	db.QueryRowContext(ctx, "Update * FROM hoge where id = ?", test) // ERROR "Use ExecContext instead of QueryRowContext to execute `UPDATE` query"
}
