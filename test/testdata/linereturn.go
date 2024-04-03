//golangcitest:args -Elinereturn
package testdata

type a struct {
	a string
}

func fooa() int {
	o := []int{0, 1}
	return o[0] // want "no blank line before"
}

func foob(s string) *a {
	o := &a{
		a: s,
	}
	return o
}

func fooc() int {
	defer fooa()

	return 0
}

func food(s string) interface{} {
	o := foob(
		s,
	)
	return o
}

func fooe() interface{} {
	o := food(
		"a",
	)
	switch s := o.(type) {
	case *a:
		return s
	default:
	}
	return o
}
