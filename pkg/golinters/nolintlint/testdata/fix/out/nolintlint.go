//golangcitest:args -Enolintlint -Elll
//golangcitest:expected_linter nolintlint
package p

import "fmt"

type FooStruct struct {
	Name string

	Value string

	// ID is an identifier.

	ID string
}

func _() {
	fmt.Println() //nolint:bob // leading space should be dropped
	fmt.Println() //nolint:bob // leading spaces should be dropped

	fmt.Println()
	fmt.Println()

	fmt.Println() //nolint:alice,lll // we don't drop individual linters from lists
}

func _() {
	//nolint:bob // leading space should be dropped
	fmt.Println()
	//nolint:bob // leading spaces should be dropped
	fmt.Println()

	fmt.Println()
	// This is a comment.

	fmt.Println()

	//nolint:alice,lll // we don't drop individual linters from lists
	fmt.Println()
}

// Foo is foo.
func Foo() {}

// Bar is bar.
//
// text text text.
func Bar() {}

func _() {}
