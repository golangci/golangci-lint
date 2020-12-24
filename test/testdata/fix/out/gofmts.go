//args: -Egofmts
package p

//gofmts:sql
const unformmattedSql = `
				SELECT
				  *
				FROM
				  mytable
				`

//gofmts:go
const unformattedGo = `const x = 1`

//gofmts:json
const unformattedJson = `
				{
				  "a": 1
				}
				`

//gofmts:sort
const unsortedValueA = 2
const unsortedValueB = 1

var unsortedSlice = []string{
	//gofmts:sort
	"a",
	"b",
}
