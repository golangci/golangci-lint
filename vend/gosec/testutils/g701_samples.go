package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG701 - SQL injection via taint analysis
var SampleCodeG701 = []CodeSample{
	{[]string{`
package main

import (
	"database/sql"
	"net/http"
)

func handler(db *sql.DB, r *http.Request) {
	name := r.URL.Query().Get("name")
	query := "SELECT * FROM users WHERE name = '" + name + "'"
	db.Query(query)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
	"net/http"
	"fmt"
)

func handler(db *sql.DB, r *http.Request) {
	id := r.FormValue("id")
	query := fmt.Sprintf("DELETE FROM users WHERE id = %s", id)
	db.Exec(query)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
)

func safeQuery(db *sql.DB) {
	// Safe - no user input
	db.Query("SELECT * FROM users")
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
	"net/http"
)

func preparedStatement(db *sql.DB, r *http.Request) {
	// Safe - using prepared statement
	name := r.URL.Query().Get("name")
	db.Query("SELECT * FROM users WHERE name = ?", name)
}
`}, 0, gosec.NewConfig()},
}
