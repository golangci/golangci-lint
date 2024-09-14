//golangcitest:args -Erecvcheck
package testdata

import "fmt"

type Bar struct{} // want `the methods of "Bar" use pointer receiver and non-pointer receiver.`

func (b Bar) A() {
	fmt.Println("A")
}

func (b *Bar) B() {
	fmt.Println("B")
}
