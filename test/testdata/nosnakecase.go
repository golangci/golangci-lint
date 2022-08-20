//golangcitest:args -Enosnakecase --internal-cmd-test
package testdata

import (
	_ "fmt"
	f_m_t "fmt" // want "f_m_t contains underscore. You should use mixedCap or MixedCap."
)

// global variable name with underscore.
var v_v = 0 // want "v_v contains underscore. You should use mixedCap or MixedCap."

// global constant name with underscore.
const c_c = 0 // want "c_c contains underscore. You should use mixedCap or MixedCap."

// struct name with underscore.
type S_a struct { // want "S_a contains underscore. You should use mixedCap or MixedCap."
	fi int
}

// non-exported struct field name with underscore.
type Sa struct {
	fi_a int // // want "fi_a contains underscore. You should use mixedCap or MixedCap."
}

// function as struct field, with parameter name with underscore.
type Sb struct {
	fib func(p_a int) // want "p_a contains underscore. You should use mixedCap or MixedCap."
}

// exported struct field with underscore.
type Sc struct {
	Fi_A int // want "Fi_A contains underscore. You should use mixedCap or MixedCap."
}

// function as struct field, with return name with underscore.
type Sd struct {
	fib func(p int) (r_a int) // want "r_a contains underscore. You should use mixedCap or MixedCap."
}

// interface name with underscore.
type I_a interface { // want "I_a contains underscore. You should use mixedCap or MixedCap."
	fn(p int)
}

// interface with parameter name with underscore.
type Ia interface {
	fn(p_a int) // want "p_a contains underscore. You should use mixedCap or MixedCap."
}

// interface with parameter name with underscore.
type Ib interface {
	Fn(p_a int) // want "p_a contains underscore. You should use mixedCap or MixedCap."
}

// function as struct field, with return name with underscore.
type Ic interface {
	Fn_a() // want "Fn_a contains underscore. You should use mixedCap or MixedCap."
}

// interface with return name with underscore.
type Id interface {
	Fn() (r_a int) // want "r_a contains underscore. You should use mixedCap or MixedCap."
}

// function name with underscore.
func f_a() {} // want "f_a contains underscore. You should use mixedCap or MixedCap."

// function's parameter name with underscore.
func fb(p_a int) {} // want "p_a contains underscore. You should use mixedCap or MixedCap."

// named return with underscore.
func fc() (r_b int) { // want "r_b contains underscore. You should use mixedCap or MixedCap."
	return 0
}

// local variable (short declaration) with underscore.
func fd(p int) int {
	v_b := p * 2 // want "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b // want "v_b contains underscore. You should use mixedCap or MixedCap."
}

// local constant with underscore.
func fe(p int) int {
	const v_b = 2 // want "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b * p // want "v_b contains underscore. You should use mixedCap or MixedCap."
}

// local variable with underscore.
func ff(p int) int {
	var v_b = 2 // want "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b * p // want "v_b contains underscore. You should use mixedCap or MixedCap."
}

// inner function, parameter name with underscore.
func fg() {
	fgl := func(p_a int) {} // want "p_a contains underscore. You should use mixedCap or MixedCap."
	fgl(1)
}

type Foo struct{}

// method name with underscore.
func (f Foo) f_a() {} // want "f_a contains underscore. You should use mixedCap or MixedCap."

// method's parameter name with underscore.
func (f Foo) fb(p_a int) {} // want "p_a contains underscore. You should use mixedCap or MixedCap."

// named return with underscore.
func (f Foo) fc() (r_b int) { return 0 } // want "r_b contains underscore. You should use mixedCap or MixedCap."

// local variable (short declaration) with underscore.
func (f Foo) fd(p int) int {
	v_b := p * 2 // want "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b // want "v_b contains underscore. You should use mixedCap or MixedCap."
}

// local constant with underscore.
func (f Foo) fe(p int) int {
	const v_b = 2 // want "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b * p // want "v_b contains underscore. You should use mixedCap or MixedCap."
}

// local variable with underscore.
func (f Foo) ff(p int) int {
	var v_b = 2 // want "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b * p // want "v_b contains underscore. You should use mixedCap or MixedCap."
}

func fna(a, p_a int) {} // want "p_a contains underscore. You should use mixedCap or MixedCap."

func fna1(a string, p_a int) {} // want "p_a contains underscore. You should use mixedCap or MixedCap."

func fnb(a, b, p_a int) {} // want "p_a contains underscore. You should use mixedCap or MixedCap."

func fnb1(a, b string, p_a int) {} // want "p_a contains underscore. You should use mixedCap or MixedCap."

func fnd(
	p_a int, // want "p_a contains underscore. You should use mixedCap or MixedCap."
	p_b int, // want "p_b contains underscore. You should use mixedCap or MixedCap."
	p_c int, // want "p_c contains underscore. You should use mixedCap or MixedCap."
) {
	f_m_t.Println("") // want "f_m_t contains underscore. You should use mixedCap or MixedCap."
}
