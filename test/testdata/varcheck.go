//golangcitest:args -Evarcheck
package testdata

var v string // ERROR "`v` is unused"
