//golangcitest:args -Einterfacebloat
package testdata

import "time"

type _ interface { // ERROR "length of interface greater than 10"
	a() time.Duration
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
	var _ interface { // ERROR "length of interface greater than 10"
		a() time.Duration
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

func __() interface { // ERROR "length of interface greater than 10"
	a() time.Duration
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
