//golangcitest:args -Erecv
//golangcitest:config_path testdata/recv_disable_type_consistency.yml
package testdata

import "fmt"

type Foo struct { // want `the methods of "Foo" use different receiver names: f, fo.`
	Name string
}

func (f Foo) A()  {}
func (fo Foo) B() {}

type Bar struct{}

func (b Bar) A()  {}
func (b *Bar) B() {}

type Fuu struct{}

func (faaa Fuu) A() { // want `the receiver name "faaa" is too long.`
	fmt.Println("a")
}
