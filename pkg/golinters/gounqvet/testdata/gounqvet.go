package testdata

import (
	"database/sql"
	"fmt"
	"strconv"
)

func badQueries() {
	query := "SELECT * FROM users" // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	
	var db *sql.DB
	rows, _ := db.Query("SELECT * FROM orders WHERE status = ?", "active") // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	_ = rows
	
	// This should not trigger because it's a COUNT function
	count := "SELECT COUNT(*) FROM users"
	_ = count
	
	// Good queries (should not trigger)
	goodQuery := "SELECT id, name, email FROM users"
	_ = goodQuery
	
	fmt.Println(query)
	
	// Use strconv to satisfy std lib import requirement
	_ = strconv.Itoa(42)
}

type SQLBuilder interface {
	Select(columns ...string) SQLBuilder
	From(table string) SQLBuilder
	Where(condition string) SQLBuilder
	Query() string
}

func badSQLBuilder(builder SQLBuilder) {
	query := builder.Select("*").From("products") // want "avoid SELECT \\* in SQL builder - explicitly specify columns to prevent unnecessary data transfer and schema change issues"
	_ = query
}

func goodSQLBuilder(builder SQLBuilder) {
	// Good usage - should not trigger
	query := builder.Select("id", "name", "price").From("products")
	_ = query
}
