package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG202 - SQL query string building via string concatenation
var SampleCodeG202 = []CodeSample{
	{[]string{`
// infixed concatenation
package main

import (
	"database/sql"
	"os"
)

func main(){
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

  q := "INSERT INTO foo (name) VALUES ('" + os.Args[0] + "')"
	rows, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
	"os"
)

func main(){
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT * FROM foo WHERE name = " + os.Args[1])
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// case insensitive match
package main

import (
	"database/sql"
	"os"
)

func main(){
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select * from foo where name = " + os.Args[1])
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// context match
package main

import (
    "context"
	"database/sql"
	"os"
)

func main(){
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	rows, err := db.QueryContext(context.Background(), "select * from foo where name = " + os.Args[1])
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// DB transaction check
package main

import (
    "context"
	"database/sql"
	"os"
)

func main(){
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(context.Background(), "select * from foo where name = " + os.Args[1])
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// DB connection check
package main

import (
    "context"
	"database/sql"
	"os"
)

func main(){
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	conn, err := db.Conn(context.Background())
	if err != nil {
		panic(err)
	}
	rows, err := conn.QueryContext(context.Background(), "select * from foo where name = " + os.Args[1])
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if err := conn.Close(); err != nil {
		panic(err)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// multiple string concatenation
package main

import (
	"database/sql"
	"os"
)

func main(){
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT * FROM foo" + "WHERE name = " + os.Args[1])
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// false positive
package main

import (
	"database/sql"
)

var staticQuery = "SELECT * FROM foo WHERE age < "
func main(){
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(staticQuery + "32")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
)

const age = "32"

var staticQuery = "SELECT * FROM foo WHERE age < "

func main(){
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
			panic(err)
	}
	rows, err := db.Query(staticQuery + age)
	if err != nil {
			panic(err)
	}
	defer rows.Close()
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

const gender = "M"
`, `
package main

import (
		"database/sql"
)

const age = "32"

var staticQuery = "SELECT * FROM foo WHERE age < "

func main(){
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
				panic(err)
		}
		rows, err := db.Query("SELECT * FROM foo WHERE gender = " + gender)
		if err != nil {
				panic(err)
		}
		defer rows.Close()
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// ExecContext match
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	result, err := db.ExecContext(context.Background(), "select * from foo where name = "+os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}`}, 1, gosec.NewConfig()},
	{[]string{`
// Exec match
package main

import (
	"database/sql"
	"fmt"
	"os"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	result, err := db.Exec("select * from foo where name = " + os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
	"fmt"
)
const gender = "M"
const age = "32"

var staticQuery = "SELECT * FROM foo WHERE age < "

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	result, err := db.Exec("SELECT * FROM foo WHERE gender = " + gender)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres password=password dbname=mydb sslmode=disable")
	if err!= nil {
		panic(err)
	}
	defer db.Close()

	var username string
	fmt.Println("请输入用户名:")
	fmt.Scanln(&username)

	var query string = "SELECT * FROM users WHERE username = '" + username + "'"
	rows, err := db.Query(query)
	if err!= nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
	"os"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	query := "SELECT * FROM album WHERE id = "
	query += os.Args[0]
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"os"
)

func main() {
	query := "SELECT * FROM album WHERE id = "
	query += os.Args[0]
	fmt.Println(query)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
	"os"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	query := "SELECT * FROM album WHERE id = "
	query = query + os.Args[0] // risky reassignment concatenation
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"database/sql"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	query := "SELECT * FROM album WHERE id = "
	query = query + "42" // safe literal reassignment concatenation
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// Shadowing edge case: tainted mutation on shadowed variable - should NOT flag
// The outer 'query' is safe and passed to db.Query.
// The inner shadowed 'query' is mutated with tainted input (irrelevant).
package main

import (
	"database/sql"
	"os"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	query := "SELECT * FROM foo WHERE id = 42" // safe outer query
	{
		query := "base"                    // shadows outer query
		query += os.Args[1]                // tainted mutation on shadow - should be ignored
		_ = query                          // prevent unused warning
	}
	rows, err := db.Query(query) // uses safe outer query
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// Shadowing edge case: no mutation on shadow, safe outer - regression guard
package main

import (
	"database/sql"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	query := "SELECT * FROM foo WHERE id = 42"
	{
		query := "shadowed but unused"
		_ = query
	}
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// package-level SQL string with tainted concatenation in init()
package main

import (
	"os"
)

var query string = "SELECT * FROM foo WHERE name = "

func init() {
	query += os.Args[1]
}
`, `
package main

import (
	"database/sql"
)

func main() {
	db, _ := sql.Open("sqlite3", ":memory:")
	_, _ = db.Query(query)
}
`}, 1, gosec.NewConfig()},
}
