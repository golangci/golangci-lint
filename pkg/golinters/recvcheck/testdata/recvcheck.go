//golangcitest:args -Erecvcheck
package testdata

type Bar struct{} // want `the methods of "Bar" use pointer receiver and non-pointer receiver.`

func (b Bar) A()  {}
func (b *Bar) B() {}
