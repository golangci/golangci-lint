//golangcitest:args -Emodernize
//golangcitest:expected_linter typecheck
package testdata

fun NotCompiles() { // want "expected declaration, found.* fun"
}
