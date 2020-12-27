//args: -Enolintlint -Elll
//config: linters-settings.nolintlint.allow-leading-space=false
package p

func nolintlint() {
	run() // nolint:bob // leading space should be dropped
	run() //  nolint:bob // leading spaces should be dropped
	// note that the next lines will retain trailing whitespace when fixed
	run() //nolint // nolint should be dropped
	run() //nolint:lll // nolint should be dropped
	run() //nolint:alice,lll,bob // enabled linter should be dropped
	run() //nolint:alice,lll,bob,nolintlint // enabled linter should not be dropped when nolintlint is nolinted
}
