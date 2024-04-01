//golangcitest:args -Egofmt
//golangcitest:config_path testdata/gofmt_no_simplify.yml
package testdata

import "fmt"

func GofmtNotSimplifiedOk() {
	var x []string
	fmt.Print(x[1:len(x)])
}

func GofmtBadFormat(){  // want "^File is not `gofmt`-ed"
}
