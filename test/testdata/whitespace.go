//args: -Ewhitespace
package testdata

func UselessStart() { // ERROR "unnecessary leading newline"

	a := 1
	_ = a
}

func UselessEnd() {
	a := 1
	_ = a

} // ERROR "unnecessary trailing newline"

func CommentsShouldBeIgnored() {
	// test
	a := 1
	_ = a
}
