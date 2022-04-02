//args: -Enonamedreturnlint
package testdata

import "fmt"

type asdf struct {
	test string
}

func noParams() {
	return
}

func argl(i string, a, b int) (ret1 string, ret2 interface{}, ret3, ret4 int, ret5 asdf) { // ERROR `named return ret1 \(string\) found in function argl`
	x := "dummy"
	return fmt.Sprintf("%s", x), nil, 1, 2, asdf{}
}

func good(i string) string {
	return i
}

func myLog(format string, args ...interface{}) {
	return
}
