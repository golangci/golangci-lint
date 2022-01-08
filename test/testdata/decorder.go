// args: -Edecorder
// config_path: testdata/configs/decorder.yml
package testdata

import "math"

const (
	decoc = math.MaxInt64
	decod = 1
)

var decoa = 1
var decob = 1 // ERROR "multiple \"var\" declarations are not allowed; use parentheses instead"

type decoe int // ERROR "type must not be placed after const"

func decof() {
	const decog = 1
}

func init() {} // ERROR "init func must be the first function in file"
