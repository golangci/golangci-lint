//args: -Egostatement
package testdata

func foo() {
	go func() { // ERROR "go statement found"

	}()

	go bar() // ERROR "go statement found"
}
