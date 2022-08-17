//golangcitest:args -Efunlen
//golangcitest:config linters-settings.funlen.lines=20
//golangcitest:config linters-settings.funlen.statements=10
package testdata

func TooManyLines() { // ERROR `Function 'TooManyLines' is too long \(22 > 20\)`
	t := struct {
		A string
		B string
		C string
		D string
		E string
		F string
		G string
		H string
		I string
	}{
		`a`,
		`b`,
		`c`,
		`d`,
		`e`,
		`f`,
		`g`,
		`h`,
		`i`,
	}
	_ = t
}

func TooManyStatements() { // ERROR `Function 'TooManyStatements' has too many statements \(11 > 10\)`
	a := 1
	b := a
	c := b
	d := c
	e := d
	f := e
	g := f
	h := g
	i := h
	j := i
	_ = j
}
