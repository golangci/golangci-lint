//golangcitest:args -Etodolint
package testdata

import "fmt"

// TODO: This is not ok // want `TODO comment should be in the form TODO\(author\)`
func todoLintNotOkFunc() {
}

// TODO(author1): This is ok
func todoLintOkFunc() {
}

type todoLintStruct struct {
	A int    // @FIXME: This field comment is not ok // want `TODO comment should be in the form FIXME\(author\)`
	B string // FIXME(author2): This field comment is ok
}

func todoLintExample() {
	// TODO(timonwong): This is ok
	//

	// 🚀🚀🚀 FixMe: 你好世界 // want `TODO comment should be in the form FIXME\(author\)`
	fmt.Println("Hello")

	fmt.Println("你好，世界") // fixme: more languages // want `TODO comment should be in the form FIXME\(author\)`

	/* TODO: old C-style comment is also supported // want `TODO comment should be in the form TODO\(author\)`
	 */

	/*
	 * TODO(timonwong) This is OK
	 *
	 */
}
