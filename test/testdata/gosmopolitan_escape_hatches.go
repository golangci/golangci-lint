//golangcitest:args -Egosmopolitan
//golangcitest:config_path testdata/configs/gosmopolitan_escape_hatches.yml
package testdata

import (
	myAlias "fmt"
)

type A string
type B = string
type C struct {
	foo string
	Bar string
}

func D(fmt string) string {
	myAlias.Println(fmt, "测试")
	return myAlias.Sprintf("%s 测试", fmt) // want `string literal contains rune in Han script`
}

type X struct {
	baz string
}

func main() {
	_ = A("测试")
	_ = string(A(string("测试")))
	_ = B("测试")
	_ = C{
		foo: "测试",
		Bar: "测试",
	}
	_ = D("测试")

	_ = &X{
		baz: "测试", // want `string literal contains rune in Han script`
	}
}
