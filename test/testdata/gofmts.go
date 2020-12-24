//args: -Egofmts
package testdata

//gofmts:sql
const unformmattedSql = `select * from mytable` // ERROR "sql formatting differs"

//gofmts:go
const unformattedGo = `const x  = 1` // ERROR "go formatting differs"

//gofmts:json
const unformattedJson = `{"a":  1}` // ERROR "json formatting differs"

var unsortedArray = []string{
	//gofmts:sort
	"z", // ERROR "block is unsorted"
	"a",
}
