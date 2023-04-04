//golangcitest:args -Etypecheck
package testdata

fun NotCompiles() { // want "expected declaration, found.* fun"
}
