//golangcitest:args -Enolintlint -Elll
//golangcitest:expected_linter nolintlint
package p

//nolint:lll // nolint should be dropped

import "fmt"

type FooStruct struct {
	//nolint:all // nolint should be dropped
	Name string

	Value string //nolint:all // nolint should be dropped

	// ID is an identifier.
	//nolint:all // nolint should be dropped (bug)
	ID string
}

func _() {
	fmt.Println() // nolint:bob // leading space should be dropped
	fmt.Println() //  nolint:bob // leading spaces should be dropped

	fmt.Println() //nolint:all // nolint should be dropped
	fmt.Println() //nolint:lll // nolint should be dropped

	fmt.Println() //nolint:alice,lll // we don't drop individual linters from lists
}

func _() {
	// nolint:bob // leading space should be dropped
	fmt.Println()
	//  nolint:bob // leading spaces should be dropped
	fmt.Println()

	//nolint:all // nolint should be dropped
	fmt.Println()
	// This is a comment.
	//nolint:lll // nolint should be dropped (bug)
	fmt.Println()

	//nolint:alice,lll // we don't drop individual linters from lists
	fmt.Println()
}

// Foo is foo.
//nolint:lll // nolint should be dropped
func Foo() {}

// Bar is bar.
//
// text text text.
//
//nolint:lll // nolint should be dropped
func Bar() {}

//nolint:lll // nolint should be dropped
func _() {}

//nolint:lll // nolint should be dropped
