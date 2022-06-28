//args: -Enosnakecase
package testdata

import (
	_ "fmt"
	f_m_t "fmt" // ERROR "f_m_t contains underscore. You should use mixedCap or MixedCap."
)

// global variable name with underscore.
var v_v = 0 // ERROR "v_v contains underscore. You should use mixedCap or MixedCap."

// global constant name with underscore.
const c_c = 0 // ERROR "c_c contains underscore. You should use mixedCap or MixedCap."

// struct name with underscore.
type S_a struct { // ERROR "S_a contains underscore. You should use mixedCap or MixedCap."
	fi int
}

// non-exported struct field name with underscore.
type Sa struct {
	fi_a int // // ERROR "fi_a contains underscore. You should use mixedCap or MixedCap."
}

// function as struct field, with parameter name with underscore.
type Sb struct {
	fib func(p_a int) // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."
}

// exported struct field with underscore.
type Sc struct {
	Fi_A int // ERROR "Fi_A contains underscore. You should use mixedCap or MixedCap."
}

// function as struct field, with return name with underscore.
type Sd struct {
	fib func(p int) (r_a int) // ERROR "r_a contains underscore. You should use mixedCap or MixedCap."
}

// interface name with underscore.
type I_a interface { // ERROR "I_a contains underscore. You should use mixedCap or MixedCap."
	fn(p int)
}

// interface with parameter name with underscore.
type Ia interface {
	fn(p_a int) // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."
}

// interface with parameter name with underscore.
type Ib interface {
	Fn(p_a int) // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."
}

// function as struct field, with return name with underscore.
type Ic interface {
	Fn_a() // ERROR "Fn_a contains underscore. You should use mixedCap or MixedCap."
}

// interface with return name with underscore.
type Id interface {
	Fn() (r_a int) // ERROR "r_a contains underscore. You should use mixedCap or MixedCap."
}

// function name with underscore.
func f_a() {} // ERROR "f_a contains underscore. You should use mixedCap or MixedCap."

// function's parameter name with underscore.
func fb(p_a int) {} // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."

// named return with underscore.
func fc() (r_b int) { // ERROR "r_b contains underscore. You should use mixedCap or MixedCap."
	return 0
}

// local variable (short declaration) with underscore.
func fd(p int) int {
	v_b := p * 2 // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."
}

// local constant with underscore.
func fe(p int) int {
	const v_b = 2 // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b * p // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."
}

// local variable with underscore.
func ff(p int) int {
	var v_b = 2 // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b * p // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."
}

// inner function, parameter name with underscore.
func fg() {
	fgl := func(p_a int) {} // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."
	fgl(1)
}

type Foo struct{}

// method name with underscore.
func (f Foo) f_a() {} // ERROR "f_a contains underscore. You should use mixedCap or MixedCap."

// method's parameter name with underscore.
func (f Foo) fb(p_a int) {} // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."

// named return with underscore.
func (f Foo) fc() (r_b int) { return 0 } // ERROR "r_b contains underscore. You should use mixedCap or MixedCap."

// local variable (short declaration) with underscore.
func (f Foo) fd(p int) int {
	v_b := p * 2 // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."
}

// local constant with underscore.
func (f Foo) fe(p int) int {
	const v_b = 2 // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b * p // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."
}

// local variable with underscore.
func (f Foo) ff(p int) int {
	var v_b = 2 // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."

	return v_b * p // ERROR "v_b contains underscore. You should use mixedCap or MixedCap."
}

func fna(a, p_a int) {} // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."

func fna1(a string, p_a int) {} // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."

func fnb(a, b, p_a int) {} // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."

func fnb1(a, b string, p_a int) {} // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."

func fnd(
	p_a int, // ERROR "p_a contains underscore. You should use mixedCap or MixedCap."
	p_b int, // ERROR "p_b contains underscore. You should use mixedCap or MixedCap."
	p_c int, // ERROR "p_c contains underscore. You should use mixedCap or MixedCap."
) {
	f_m_t.Println("") // ERROR "f_m_t contains underscore. You should use mixedCap or MixedCap."
}
