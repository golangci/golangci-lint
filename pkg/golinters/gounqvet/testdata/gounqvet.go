package testdata

import (
	"database/sql"
	"fmt"
	"strconv"
)

func badQueries() {
	// want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	query := "SELECT * FROM users"
	
	var db *sql.DB
	// want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	rows, _ := db.Query("SELECT * FROM orders WHERE status = ?", "active")
	_ = rows
	
	// This should not trigger because it's a COUNT function
	count := "SELECT COUNT(*) FROM users"
	_ = count
	
	// This should not trigger because of nolint comment
	debug := "SELECT * FROM debug_table" //nolint:gounqvet
	_ = debug
	
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
	// want "avoid SELECT \\* in SQL builder - explicitly specify columns to prevent unnecessary data transfer and schema change issues"
	query := builder.Select("*").From("products")
	_ = query
}

func goodSQLBuilder(builder SQLBuilder) {
	// Good usage - should not trigger
	query := builder.Select("id", "name", "price").From("products")
	_ = query
}