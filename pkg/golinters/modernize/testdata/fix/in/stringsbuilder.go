//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package stringsbuilder

// basic test
func _() {
	var s string
	s += "before"
	for range 10 {
		s += "in" // want "using string \\+= string in a loop is inefficient"
		s += "in2"
	}
	s += "after"
	print(s)
}

// with initializer
func _() {
	var s = "a"
	for range 10 {
		s += "b" // want "using string \\+= string in a loop is inefficient"
	}
	print(s)
}

// with empty initializer
func _() {
	var s = ""
	for range 10 {
		s += "b" // want "using string \\+= string in a loop is inefficient"
	}
	print(s)
}

// with short decl
func _() {
	s := "a"
	for range 10 {
		s += "b" // want "using string \\+= string in a loop is inefficient"
	}
	print(s)
}

// with short decl and empty initializer
func _() {
	s := ""
	for range 10 {
		s += "b" // want "using string \\+= string in a loop is inefficient"
	}
	print(s)
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
		s := "a"
		for range 10 {
			s += "b" // want "using string \\+= string in a loop is inefficient"
		}
		print(s)
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
		msg string
	)
	for _, s := range slice {
		msg += s // want "using string \\+= string in a loop is inefficient"
	}
	return msg
}
