//golangcitest:args -Eunqueryvet
//golangcitest:config_path testdata/unqueryvet_custom.yml
package testdata

import (
	"database/sql"
	"fmt"
	"strconv"
)

func _() {
	query := "SELECT * FROM users" // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"

	var db *sql.DB
	rows, _ := db.Query("SELECT * FROM orders WHERE status = ?", "active") // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	_ = rows

	count := "SELECT COUNT(*) FROM users"
	_ = count

	goodQuery := "SELECT id, name, email FROM users"
	_ = goodQuery

	fmt.Println(query)

	_ = strconv.Itoa(42)
}

// Custom allowed patterns test - SELECT * from temp tables should be allowed
