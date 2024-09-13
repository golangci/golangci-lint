//golangcitest:args -Erecv
//golangcitest:config_path testdata/recv_max_name_length.yml
package testdata

import "fmt"

type Fuu struct{}

func (faaa Fuu) A() {
	fmt.Println("a")
}

type Foo struct{}

func (faaaaaa Foo) A() { // want `the receiver name "faaaaaa" is too long.`
	fmt.Println("a")
}
