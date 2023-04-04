//golangcitest:args -Evarcheck --internal-cmd-test
package testdata

var v string // want "`v` is unused"
