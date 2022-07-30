//golangcitest:args -Einterfacebloat
package testdata

type _ interface { // want "length of interface greater than 10"
	a()
	b()
	c()
	d()
	f()
	g()
	h()
	i()
	j()
	k()
	l()
}

func _() {
	var _ interface { // want "length of interface greater than 10"
		a()
		b()
		c()
		d()
		f()
		g()
		h()
		i()
		j()
		k()
		l()
	}
}

func __() interface { // want "length of interface greater than 10"
	a()
	b()
	c()
	d()
	f()
	g()
	h()
	i()
	j()
	k()
	l()
} {
	return nil
}
