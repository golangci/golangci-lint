//args: -Egofmts
package p

//gofmts:sql
const unformmattedSql = `select * from mytable`

//gofmts:go
const unformattedGo = `const x  = 1`

//gofmts:json
const unformattedJson = `{"a":  1}`

//gofmts:sort
const unsortedValueB = 1
const unsortedValueA = 2

var unsortedSlice = []string{
	//gofmts:sort
	"b",
	"a",
}
