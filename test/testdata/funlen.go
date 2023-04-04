//golangcitest:args -Efunlen
//golangcitest:config_path testdata/configs/funlen.yml
package testdata

func TooManyLines() { // want `Function 'TooManyLines' is too long \(22 > 20\)`
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

func TooManyStatements() { // want `Function 'TooManyStatements' has too many statements \(11 > 10\)`
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
