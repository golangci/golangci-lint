//golangcitest:args -Ewsl_v5
//golangcitest:expected_exitcode 0
package testdata

import (
	"bytes"
	"encoding/json"
)

type T struct {
	I int
}

func NewT() *T {
	return &T{}
}

func (*T) Fn() int {
	return 1
}

func FnC(_ string, fn func() error) error {
	fn()
	return nil
}

func strictAppend() {
	s := []string{}
	s = append(s, "a")
	s = append(s, "b")
	x := "c"
	s = append(s, x)
	y := "e"

	s = append(s, "d") // want `missing whitespace above this line \(no shared variables above append\)`
	s = append(s, y)
}

func incDec() {
	x := 1
	x++
	x--
	y := x

	_ = y
}

func assignAndCall() {
	t1 := NewT()
	t2 := NewT()

	t1.Fn()
	x := t1.Fn()
	t1.Fn()

	_, _ = t2, x
}

func closureInCall() {
	buf := &bytes.Buffer{}
	_ = FnC("buf", func() error {
		return json.NewEncoder(buf).Encode("x")
	})
}

func assignAfterBlock() {
	x := 1
	if x > 0 {
		return
	}

	x = 2 // want `missing whitespace above this line \(invalid statement above assign\)`
}

func decl() {
	var x string

	y := "" // want `missing whitespace above this line \(invalid statement above assign\)`

	_, _ = x, y
}
