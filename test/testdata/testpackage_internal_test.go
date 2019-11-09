//args: -Etestpackage -Egochecknoglobals
package testdata

// Test expects at least one issue in the file.
// So we have to add global variable and enable gochecknoglobals.
var global = `global` // ERROR "`global` is a global variable"
