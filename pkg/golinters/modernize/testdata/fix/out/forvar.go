//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package forvar

func _(m map[int]int, s []int) {
	// changed
	for i := range s {
		go f(i)
	}
	for _, v := range s {
		go f(v)
	}
	for k, v := range m {
		go f(k)
		go f(v)
	}
	for k, v := range m {
		go f(k)
		go f(v)
	}
	for k, v := range m {
		go f(k)
		go f(v)
	}
	for k, v := range m {
		go f(k)
		go f(v)
	}
	for i := range s {
		go f(i)
	}
	// nope
	var i, k, v int

	for i = range s { // nope, scope change
		i := i
		go f(i)
	}
	for _, v = range s { // nope, scope change
		v := v
		go f(v)
	}
	for k = range m { // nope, scope change
		k := k
		go f(k)
	}
	for k, v = range m { // nope, scope change
		k := k
		v := v
		go f(k)
		go f(v)
	}
	for _, v = range m { // nope, scope change
		v := v
		go f(v)
	}
	for _, v = range m { // nope, not x := x
		v := i
		go f(v)
	}
	for k, v := range m { // nope, LHS and RHS differ
		v, k := k, v
		go f(k)
		go f(v)
	}
	for k, v := range m { // nope, not a simple redecl
		k, v, x := k, v, 1
		go f(k)
		go f(v)
		go f(x)
	}
	for i := range s { // nope, not a simple redecl
		i := (i)
		go f(i)
	}
	for i := range s { // nope, not a simple redecl
		i := i + 1
		go f(i)
	}
}

func f(n int) {}
