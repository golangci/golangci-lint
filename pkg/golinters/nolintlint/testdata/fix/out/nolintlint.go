//golangcitest:args -Enolintlint -Elll
//golangcitest:expected_linter nolintlint
package p

import "fmt"

func nolintlint() {
	fmt.Println() //nolint:bob // leading space should be dropped
	fmt.Println() //nolint:bob // leading spaces should be dropped

	fmt.Println()
	fmt.Println()

	fmt.Println() //nolint:alice,lll // we don't drop individual linters from lists
}
