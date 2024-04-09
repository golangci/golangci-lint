//golangcitest:args -Edecorder
//golangcitest:config_path testdata/decorder_custom.yml
package testdata

import "math"

const (
	decoc = math.MaxInt64
	decod = 1
)

var decoa = 1
var decob = 1 // want "multiple \"var\" declarations are not allowed; use parentheses instead"

type decoe int // want "type must not be placed after const"

func decof() {
	const decog = 1
}

func init() {} // want "init func must be the first function in file"
