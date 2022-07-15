//golangcitest:args -Esqlclosecheck
package testdata

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

var (
	ctx    context.Context
	db     *sql.DB
	age    = 27
	userID = 43
)

func rowsCorrectDeferBlock() {

	rows, err := db.QueryContext(ctx, "SELECT name FROM users WHERE age=?", age)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Print("problem closing rows")
		}
	}()

	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s are %d years old", strings.Join(names, ", "), age)
}

func rowsCorrectDefer() {
	rows, err := db.QueryContext(ctx, "SELECT name FROM users WHERE age=?", age)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s are %d years old", strings.Join(names, ", "), age)
}

func rowsMissingClose() {
	rows, err := db.QueryContext(ctx, "SELECT name FROM users WHERE age=?", age) // ERROR "Rows/Stmt was not closed"
	if err != nil {
		log.Fatal(err)
	}
	// defer rows.Close()

	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s are %d years old", strings.Join(names, ", "), age)
}

func rowsNonDeferClose() {
	rows, err := db.QueryContext(ctx, "SELECT name FROM users WHERE age=?", age)
	if err != nil {
		log.Fatal(err)
	}

	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s are %d years old", strings.Join(names, ", "), age)

	rows.Close() // ERROR "Close should use defer"
}

func rowsPassedAndClosed() {
	rows, err := db.QueryContext(ctx, "SELECT name FROM users")
	if err != nil {
		log.Fatal(err)
	}

	rowsClosedPassed(rows)
}

func rowsClosedPassed(rows *sql.Rows) {
	rows.Close()
}

func rowsPassedAndNotClosed(rows *sql.Rows) {
	rows, err := db.QueryContext(ctx, "SELECT name FROM users")
	if err != nil {
		log.Fatal(err)
	}

	rowsDontClosedPassed(rows)
}

func rowsDontClosedPassed(*sql.Rows) {

}

func rowsReturn() (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx, "SELECT name FROM users WHERE age=?", age)
	if err != nil {
		log.Fatal(err)
	}
	return rows, nil
}

func rowsReturnShort() (*sql.Rows, error) {
	return db.QueryContext(ctx, "SELECT name FROM users WHERE age=?", age)
}

func stmtCorrectDeferBlock() {
	// In normal use, create one Stmt when your process starts.
	stmt, err := db.PrepareContext(ctx, "SELECT username FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Print("problem closing stmt")
		}
	}()

	// Then reuse it each time you need to issue the query.
	var username string
	err = stmt.QueryRowContext(ctx, userID).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		log.Fatalf("no user with id %d", userID)
	case err != nil:
		log.Fatal(err)
	default:
		log.Printf("username is %s\n", username)
	}
}

func stmtCorrectDefer() {
	// In normal use, create one Stmt when your process starts.
	stmt, err := db.PrepareContext(ctx, "SELECT username FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Then reuse it each time you need to issue the query.
	var username string
	err = stmt.QueryRowContext(ctx, userID).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		log.Fatalf("no user with id %d", userID)
	case err != nil:
		log.Fatal(err)
	default:
		log.Printf("username is %s\n", username)
	}
}

func stmtMissingClose() {
	// In normal use, create one Stmt when your process starts.
	stmt, err := db.PrepareContext(ctx, "SELECT username FROM users WHERE id = ?") // ERROR "Rows/Stmt was not closed"
	if err != nil {
		log.Fatal(err)
	}
	// defer stmt.Close()

	// Then reuse it each time you need to issue the query.
	var username string
	err = stmt.QueryRowContext(ctx, userID).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		log.Fatalf("no user with id %d", userID)
	case err != nil:
		log.Fatal(err)
	default:
		log.Printf("username is %s\n", username)
	}
}

func stmtNonDeferClose() {
	// In normal use, create one Stmt when your process starts.
	stmt, err := db.PrepareContext(ctx, "SELECT username FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	// Then reuse it each time you need to issue the query.
	var username string
	err = stmt.QueryRowContext(ctx, userID).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		log.Fatalf("no user with id %d", userID)
	case err != nil:
		log.Fatal(err)
	default:
		log.Printf("username is %s\n", username)
	}

	stmt.Close() // ERROR "Close should use defer"
}

func stmtReturn() (*sql.Stmt, error) {
	stmt, err := db.PrepareContext(ctx, "SELECT username FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func stmtReturnShort() (*sql.Stmt, error) {
	return db.PrepareContext(ctx, "SELECT username FROM users WHERE id = ?")
}
