//golangcitest:args -Egofmt
//golangcitest:config linters-settings.gofmt.simplify=false
package testdata

import "fmt"

func GofmtNotSimplifiedOk() {
	var x []string
	fmt.Print(x[1:len(x)])
}

func GofmtBadFormat(){  // ERROR "^File is not `gofmt`-ed"
}
