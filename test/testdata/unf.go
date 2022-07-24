//golangcitest:args -Eunf
package testdata

import "fmt"

func printf() {
	fmt.Printf("") // ERROR "format like function Printf used without arguments"
}

func errorf() error {
	return fmt.Errorf("") // ERROR "format like function Errorf used without arguments"
}

type s struct{}

func (*s) msgf(format string, a ...string) {
	fmt.Printf(format, a)
}

func _() {
	var s s
	s.msgf("") // ERROR "format like function msgf used without arguments"
}
