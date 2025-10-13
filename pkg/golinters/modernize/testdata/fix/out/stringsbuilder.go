//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package stringsbuilder

import "strings"

// basic test
func _() {
	var s strings.Builder
	s.WriteString("before")
	for range 10 {
		s.WriteString("in") // want "using string \\+= string in a loop is inefficient"
		s.WriteString("in2")
	}
	s.WriteString("after")
	print(s.String())
}

// with initializer
func _() {
	var s strings.Builder
	s.WriteString("a")
	for range 10 {
		s.WriteString("b") // want "using string \\+= string in a loop is inefficient"
	}
	print(s.String())
}

// with empty initializer
func _() {
	var s strings.Builder
	for range 10 {
		s.WriteString("b") // want "using string \\+= string in a loop is inefficient"
	}
	print(s.String())
}

// with short decl
func _() {
	var s strings.Builder
	s.WriteString("a")
	for range 10 {
		s.WriteString("b") // want "using string \\+= string in a loop is inefficient"
	}
	print(s.String())
}

// with short decl and empty initializer
func _() {
	var s strings.Builder
	for range 10 {
		s.WriteString("b") // want "using string \\+= string in a loop is inefficient"
	}
	print(s.String())
}

// nope: += must appear at least once within a loop.
func _() {
	var s string
	s += "a"
	s += "b"
	s += "c"
	print(s)
}

// nope: the declaration of s is not in a block.
func _() {
	if s := "a"; true {
		for range 10 {
			s += "x"
		}
		print(s)
	}
}

// in a switch (special case of "in a block" logic)
func _() {
	switch {
	default:
		var s strings.Builder
		s.WriteString("a")
		for range 10 {
			s.WriteString("b") // want "using string \\+= string in a loop is inefficient"
		}
		print(s.String())
	}
}

// nope: don't handle direct assignments to the string  (only +=).
func _(x string) string {
	var s string
	s = x
	for range 3 {
		s += "" + x
	}
	return s
}

// Regression test for bug in a GenDecl with parens.
func issue75318(slice []string) string {
	var (
		msg strings.Builder
	)
	for _, s := range slice {
		msg.WriteString(s) // want "using string \\+= string in a loop is inefficient"
	}
	return msg.String()
}
