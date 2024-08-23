//golangcitest:args -Einterfaceguard
package testdata

import "fmt"

type Greeter interface {
	SayHi() string
}

type BasicGreeter struct{}

func (b *BasicGreeter) SayHi() string {
	return "Welcome!"
}

var (
	greeter      Greeter
	otherGreeter Greeter
)

func nilBasicComparison() {
	if greeter == nil { // want "comparing interface to nil"
		fmt.Println("greeter is nil")
	}

	if greeter != nil { // want "comparing interface to nil"
		fmt.Println("greeter is not nil")
	}
}

func interfaceBasicComparison() {
	if greeter == otherGreeter { // want "comparing two interfaces"
		fmt.Println("greeter == otherGreeter")
	}

	if greeter != otherGreeter { // want "comparing two interfaces"
		fmt.Println("greeter != otherGreeter")
	}
}
