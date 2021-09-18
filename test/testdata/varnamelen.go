//args: -Evarnamelen
package testdata

func varnamelen() {
	x := 1 // ERROR "variable name 'x' is too short for the scope of its usage"
	x++
	x++
	x++
	x++
	x++
	x++
	x++
	x++
	x++
}
