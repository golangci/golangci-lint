//golangcitest:args -Eunqueryvet
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

type SQLBuilder interface {
	Select(columns ...string) SQLBuilder
	From(table string) SQLBuilder
	Where(condition string) SQLBuilder
	Query() string
}

func _(builder SQLBuilder) {
	query := builder.Select("*").From("products") // want "avoid SELECT \\* in SQL builder - explicitly specify columns to prevent unnecessary data transfer and schema change issues"
	_ = query
}

func _(builder SQLBuilder) {
	query := builder.Select("id", "name", "price").From("products")
	_ = query
}
